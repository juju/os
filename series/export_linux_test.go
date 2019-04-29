// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series

var (
	DistroInfo    = &distroInfo
	ReadSeries    = readSeries
	OSReleaseFile = &osReleaseFile
)

func HideUbuntuSeries() func() {
	origSeries := ubuntuSeries
	ubuntuSeries = make(map[string]seriesVersion)
	return func() {
		ubuntuSeries = origSeries
	}
}
