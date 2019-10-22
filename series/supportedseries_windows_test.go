// Copyright 2014 Canonical Ltd.
// Copyright 2014 Cloudbase Solutions SRL
// Licensed under the LGPLv3, see LICENCE file for details.

package series_test

import (
	"github.com/juju/os/series"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type supportedSeriesWindowsSuite struct {
}

var _ = gc.Suite(&supportedSeriesWindowsSuite{})

func (s *supportedSeriesWindowsSuite) TestSeriesVersion(c *gc.C) {
	vers, err := series.SeriesVersion("win8")
	if err != nil {
		c.Assert(err, gc.Not(gc.ErrorMatches), `invalid series "win8"`, gc.Commentf(`unable to lookup series "win8"`))
	} else {
		c.Assert(err, jc.ErrorIsNil)
	}
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(vers, gc.Equals, "win8")
}

func (s *supportedSeriesSuite) TestIsWindowsNano(c *gc.C) {
	var isWindowsNanoTests = []struct {
		series   string
		expected bool
	}{
		{"win2016nano", true},
		{"win2016", false},
		{"win2012r2", false},
		{"trusty", false},
	}

	for _, t := range isWindowsNanoTests {
		c.Assert(series.IsWindowsNano(t.series), gc.Equals, t.expected)
	}
}

func (s supportedSeriesWindowsSuite) TestWindowsVersions(c *gc.C) {
	windowsVersions, overwrittenValues := series.WindowsVersions()
	wlen := len(series.WindowsVersionMap)
	nlen := len(series.WindowsNanoMap)
	verify := 0

	// This return len(series.WindowsVersionMap) + n,
	// n equals to the number of values we have overwritten
	// because we overwrite a value of the map in series.WindowsVersions
	for i, ival := range windowsVersions {
		for j, jval := range series.WindowsVersionMap {
			if i == j && ival == jval {
				verify++
			}
		}
	}
	c.Assert(verify+len(overwrittenValues), gc.Equals, wlen)

	verify = 0
	// This should return len(WindowsNanoMap)
	for i, ival := range windowsVersions {
		for j, jval := range series.WindowsNanoMap {
			if i == j && ival == jval {
				verify++
			}
		}
	}
	c.Assert(verify, gc.Equals, nlen)
}
