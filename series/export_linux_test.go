// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series

var (
	UbuntuDistroInfoPath = &UbuntuDistroInfo
	ReadSeries           = readSeries
	OSReleaseFile        = &osReleaseFile
)

// HideUbuntuSeries hides the global state of the ubuntu series for tests. The
// function returns a closure, that puts the global state back once called.
// This is not concurrent safe.
func HideUbuntuSeries() func() {
	origSeries := ubuntuSeries
	ubuntuSeries = make(map[string]SeriesVersionInfo)
	return func() {
		ubuntuSeries = origSeries
	}
}
