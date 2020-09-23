package series

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/os"
)

type System struct {
	OS       os.OSType
	Version  string
	Resource string
}

// String respresentation of the System, used for series backwards compatability.
func (s *System) String() string {
	// Handle legacy series.
	if s.OS != os.Unknown && s.Version != "" && s.Resource == "" {
		// Try matching version to a series.
		if seriesFromVersion, err := VersionSeries(s.Version); err == nil {
			if osFromSeries, err := GetOSFromSeries(seriesFromVersion); err == nil && osFromSeries == s.OS {
				return seriesFromVersion
			}
		}
		// Try matching version as a series.
		if osFromSeries, err := GetOSFromSeries(s.Version); err == nil && osFromSeries == s.OS {
			return s.Version
		}
	}
	// Handle new system as a series.
	str := "system"
	if s.OS != os.Unknown {
		str += fmt.Sprintf("#os=%s", s.OS.SystemString())
	}
	if s.Version != "" {
		str += fmt.Sprintf("#version=%s", s.Version)
	}
	if s.Resource != "" {
		str += fmt.Sprintf("#resource=%s", s.Resource)
	}
	return str
}

// regex to match k=v from system series strings in the form
// "system#os=ubuntu#version=18.04#resource=imagename"
var systemSeriesRegex = regexp.MustCompile(`#(os|version|resource)=([^#]+)`)

// ParseSystemFromSeries matches legacy series like "focal" or parses a system as series string
// in the form "system#os=ubuntu#version=18.04#resource=imagename"
func ParseSystemFromSeries(s string) (System, error) {
	var err error
	if !strings.HasPrefix(s, "system") {
		osType, err := GetOSFromSeries(s)
		if err != nil {
			return System{}, errors.Trace(err)
		}
		version, err := SeriesVersion(s)
		if err != nil {
			return System{}, errors.Trace(err)
		}
		return System{
			OS:      osType,
			Version: version,
		}, nil
	}
	propString := strings.TrimPrefix(s, "system")
	matches := systemSeriesRegex.FindAllStringSubmatch(propString, -1)
	if len(matches) == 0 {
		return System{}, errors.NotValidf("invalid system series string %q", s)
	}
	matchedCharacters := 0
	system := System{}
	for _, v := range matches {
		matchedCharacters += len(v[0])
		key := v[1]
		value := v[2]
		switch key {
		case "os":
			system.OS, err = os.ParseSystemOS(value)
			if err != nil {
				return System{}, errors.Trace(err)
			}
		case "version":
			system.Version = value
		case "resource":
			system.Resource = value
		default:
			return System{}, errors.NotValidf("invalid key %q in system series string %q", key, s)
		}
	}
	if matchedCharacters != len(propString) {
		return System{}, errors.NotValidf("invalid system series string %q", s)
	}
	return system, nil
}
