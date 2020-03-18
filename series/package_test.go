// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series_test

import (
	"testing"

	gc "gopkg.in/check.v1"
)

//go:generate mockgen -package series -destination filesystem_mock_test.go github.com/juju/os/series FileSystem

func Test(t *testing.T) {
	gc.TestingT(t)
}
