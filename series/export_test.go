// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series

var (
	KernelToMajor                 = kernelToMajor
	MacOSXSeriesFromKernelVersion = macOSXSeriesFromKernelVersion
	MacOSXSeriesFromMajorVersion  = macOSXSeriesFromMajorVersion

	origSeriesVersions map[string]string
	origUbuntuSeries   map[string]seriesVersion
)

func init() {
	origSeriesVersions = make(map[string]string)
	for k, v := range seriesVersions {
		origSeriesVersions[k] = v
	}
	origUbuntuSeries = make(map[string]seriesVersion)
	for k, v := range ubuntuSeries {
		origUbuntuSeries[k] = v
	}
}

func SetSeriesVersions(value map[string]string) func() {
	origVersions := seriesVersions
	origUpdated := updatedseriesVersions
	seriesVersions = value
	updateVersionSeries()
	updatedseriesVersions = len(value) != 0
	return func() {
		seriesVersions = origVersions
		updateVersionSeries()
		updatedseriesVersions = origUpdated
	}
}

// RestoreSeriesVersions restores the value of seriesVersions to a copy of the
// original, for use in test cleanup.
func RestoreSeriesVersions() {
	seriesVersions = make(map[string]string)
	for k, v := range origSeriesVersions {
		seriesVersions[k] = v
	}
}

// UbuntuSupportedSeries exports the ubuntuSeries for testing.
func UbuntuSupportedSeries() map[string]seriesVersion {
	return ubuntuSeries
}

// RestoreUbuntuSeries restore the value of ubuntuSeries to a copy of the
// original, for use in test cleanup.
func RestoreUbuntuSeries() {
	ubuntuSeries = make(map[string]seriesVersion)
	for k, v := range origUbuntuSeries {
		ubuntuSeries[k] = v
	}
}
