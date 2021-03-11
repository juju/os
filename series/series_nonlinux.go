// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

// +build !linux

package series

import (
	"os"

	jujuos "github.com/juju/os"
)

// TODO(ericsnow) Refactor dependents so we can remove this for non-linux.

// ReleaseVersion is a function that has no meaning except on linux.
func ReleaseVersion() string {
	return ""
}

// // LocalSeriesVersionInfo is a function that has no meaning except on Linux.
func LocalSeriesVersionInfo() (jujuos.OSType, map[string]SeriesVersionInfo, error) {
	return jujuos.Unknown, nil, nil
}

func updateLocalSeriesVersions() error {
	return nil
}

// defaultFileSystem implements the FileSystem for the DistroInfo.
type defaultFileSystem struct{}

func (defaultFileSystem) Open(path string) (*os.File, error) {
	return nil, os.ErrNotExist
}

func (defaultFileSystem) Exists(path string) bool {
	return false
}
