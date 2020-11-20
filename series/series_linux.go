// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/juju/errors"
	jujuos "github.com/juju/os"
)

var (
	// osReleaseFile is the name of the file that is read in order to determine
	// the linux type release version.
	osReleaseFile = "/etc/os-release"

	// timeNow is time.Now, but overrideable via TimeNow in tests.
	timeNow = time.Now
)

const (
	// this is just for an approximation in an error case, when the eol
	// date has a parse error.
	day  = 24 * time.Hour
	year = 365 * day
)

func readSeries() (string, error) {
	values, err := jujuos.ReadOSRelease(osReleaseFile)
	if err != nil {
		return "unknown", err
	}
	updateSeriesVersionsOnce()
	return seriesFromOSRelease(values)
}

func seriesFromOSRelease(values map[string]string) (string, error) {
	switch values["ID"] {
	case strings.ToLower(jujuos.Ubuntu.String()):
		return getValueFromSeriesVersion(ubuntuSeries, values["VERSION_ID"])
	case strings.ToLower(jujuos.CentOS.String()):
		codename := fmt.Sprintf("%s%s", values["ID"], values["VERSION_ID"])
		return getValue(centosSeries, codename)
	case strings.ToLower(jujuos.OpenSUSE.String()):
		codename := fmt.Sprintf("%s%s",
			values["ID"],
			strings.Split(values["VERSION_ID"], ".")[0])
		return getValue(opensuseSeries, codename)
	default:
		return genericLinuxSeries, nil
	}
}

func getValue(from map[string]string, val string) (string, error) {
	for serie, ver := range from {
		if ver == val {
			return serie, nil
		}
	}
	return "unknown", errors.New("could not determine series")
}

func getValueFromSeriesVersion(from map[string]seriesVersion, val string) (string, error) {
	for s, version := range from {
		if version.Version == val {
			return s, nil
		}
	}
	return "unknown", errors.New("could not determine series")
}

// ReleaseVersion looks for the value of VERSION_ID in the content of
// the os-release.  If the value is not found, the file is not found, or
// an error occurs reading the file, an empty string is returned.
func ReleaseVersion() string {
	release, err := jujuos.ReadOSRelease(osReleaseFile)
	if err != nil {
		return ""
	}
	return release["VERSION_ID"]
}

// updateLocalSeriesVersions updates seriesVersions from
// /usr/share/distro-info/ubuntu.csv if possible..
func updateLocalSeriesVersions() error {
	distroInfo := NewDistroInfo(UbuntuDistroInfo)
	if err := distroInfo.Refresh(); err != nil {
		return errors.Trace(err)
	}

	now := timeNow().UTC()

	for seriesName, version := range distroInfo.info {
		var esm bool
		if existing, ok := ubuntuSeries[seriesName]; ok {
			esm = existing.ESMSupported
		}

		// The numeric version may contain a LTS moniker so strip that out.
		trimmedVersion := strings.TrimSuffix(version.Version, " LTS")
		seriesVersions[seriesName] = trimmedVersion

		// If the series already exists inside of ubuntuSeries then don't
		// overwrite that existing one, except to update the supported status.
		supported := version.Supported(now)

		if us, ok := ubuntuSeries[seriesName]; ok {
			us.Supported = supported
			ubuntuSeries[seriesName] = us
			continue
		}

		ubuntuSeries[seriesName] = seriesVersion{
			Version:                  version.Version,
			Supported:                supported,
			ESMSupported:             esm,
			LTS:                      version.LTS(),
			CreatedByLocalDistroInfo: true,
		}
	}

	return nil
}

// defaultFileSystem implements the FileSystem for the DistroInfo.
type defaultFileSystem struct{}

func (defaultFileSystem) Open(path string) (*os.File, error) {
	return os.Open(path)
}

func (defaultFileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
