// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series_test

import (
	"time"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/os/v2"
	"github.com/juju/os/v2/series"
)

type supportedSeriesSuite struct {
	testing.CleanupSuite
}

var _ = gc.Suite(&supportedSeriesSuite{})

func (s *supportedSeriesSuite) SetUpTest(c *gc.C) {
	s.CleanupSuite.SetUpTest(c)

	cleanup := series.SetSeriesVersions(make(map[string]string))
	s.AddCleanup(func(*gc.C) { cleanup() })

	s.PatchValue(series.TimeNow, func() time.Time {
		return time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC)
	})
}

var getOSFromSeriesTests = []struct {
	series string
	want   os.OSType
	err    string
}{{
	series: "precise",
	want:   os.Ubuntu,
}, {
	series: "win2012r2",
	want:   os.Windows,
}, {
	series: "win2016nano",
	want:   os.Windows,
}, {
	series: "mountainlion",
	want:   os.OSX,
}, {
	series: "centos7",
	want:   os.CentOS,
}, {
	series: "opensuseleap",
	want:   os.OpenSUSE,
}, {
	series: "kubernetes",
	want:   os.Kubernetes,
}, {
	series: "genericlinux",
	want:   os.GenericLinux,
}, {
	series: "",
	err:    "series \"\" not valid",
},
}

func (s *supportedSeriesSuite) TestGetOSFromSeries(c *gc.C) {
	for _, t := range getOSFromSeriesTests {
		got, err := series.GetOSFromSeries(t.series)
		if t.err != "" {
			c.Assert(err, gc.ErrorMatches, t.err)
		} else {
			c.Check(err, jc.ErrorIsNil)
			c.Assert(got, gc.Equals, t.want)
		}
	}
}

func (s *supportedSeriesSuite) TestUnknownOSFromSeries(c *gc.C) {
	_, err := series.GetOSFromSeries("Xuanhuaceratops")
	c.Assert(err, jc.Satisfies, series.IsUnknownOSForSeriesError)
	c.Assert(err, gc.ErrorMatches, `unknown OS for series: "Xuanhuaceratops"`)
}

var getOSFromSeriesBaseOSTests = []struct {
	series, baseOS string
	want           os.OSType
	err            string
}{{
	series: "precise",
	baseOS: "ubuntu",
	want:   os.Ubuntu,
}, {
	series: "win2012r2",
	baseOS: "windows",
	want:   os.Windows,
}, {
	series: "win2016nano",
	baseOS: "windows",
	want:   os.Windows,
}, {
	series: "mountainlion",
	baseOS: "osx",
	want:   os.OSX,
}, {
	series: "7",
	baseOS: "centos",
	want:   os.CentOS,
}, {
	series: "opensuseleap",
	baseOS: "opensuse",
	want:   os.OpenSUSE,
}, {
	series: "kubernetes",
	baseOS: "kubernetes",
	want:   os.Kubernetes,
}, {
	series: "genericlinux",
	baseOS: "genericlinux",
	want:   os.GenericLinux,
}, {
	series: "",
	err:    "series \"\" not valid",
},
}

func (s *supportedSeriesSuite) TestGetOSFromSeriesWithBaseOS(c *gc.C) {
	for _, t := range getOSFromSeriesBaseOSTests {
		c.Logf("series %q os %q", t.series, t.baseOS)
		got, err := series.GetOSFromSeriesWithBaseOS(t.series, t.baseOS)
		if t.err != "" {
			c.Assert(err, gc.ErrorMatches, t.err)
		} else {
			c.Check(err, jc.ErrorIsNil)
			c.Assert(got, gc.Equals, t.want)
		}
	}
}

func setSeriesTestData() {
	series.SetSeriesVersions(map[string]string{
		"trusty":       "14.04",
		"utopic":       "14.10",
		"win7":         "win7",
		"win81":        "win81",
		"win2016nano":  "win2016nano",
		"centos7":      "centos7",
		"opensuseleap": "opensuse42",
		"genericlinux": "genericlinux",
	})
}

func (s *supportedSeriesSuite) TestOSSupportedSeries(c *gc.C) {
	setSeriesTestData()
	supported := series.OSSupportedSeries(os.Ubuntu)
	c.Assert(supported, jc.SameContents, []string{"trusty", "utopic"})
	supported = series.OSSupportedSeries(os.Windows)
	c.Assert(supported, jc.SameContents, []string{"win7", "win81", "win2016nano"})
	supported = series.OSSupportedSeries(os.CentOS)
	c.Assert(supported, jc.SameContents, []string{"centos7"})
	supported = series.OSSupportedSeries(os.OpenSUSE)
	c.Assert(supported, jc.SameContents, []string{"opensuseleap"})
	supported = series.OSSupportedSeries(os.GenericLinux)
	c.Assert(supported, jc.SameContents, []string{"genericlinux"})
}

func (s *supportedSeriesSuite) TestVersionSeriesValid(c *gc.C) {
	setSeriesTestData()
	seriesResult, err := series.VersionSeries("14.04")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert("trusty", gc.DeepEquals, seriesResult)
}

func (s *supportedSeriesSuite) TestVersionSeriesEmpty(c *gc.C) {
	setSeriesTestData()
	_, err := series.VersionSeries("")
	c.Assert(err, gc.ErrorMatches, `.*unknown series for version: "".*`)
}

func (s *supportedSeriesSuite) TestVersionSeriesInvalid(c *gc.C) {
	setSeriesTestData()
	_, err := series.VersionSeries("73655")
	c.Assert(err, gc.ErrorMatches, `.*unknown series for version: "73655".*`)
}

func (s *supportedSeriesSuite) TestSeriesVersionEmpty(c *gc.C) {
	setSeriesTestData()
	_, err := series.SeriesVersion("")
	c.Assert(err, gc.ErrorMatches, `.*unknown version for series: "".*`)
}

func (s *supportedSeriesSuite) TestUbuntuSeriesVersionEmpty(c *gc.C) {
	_, err := series.UbuntuSeriesVersion("")
	c.Assert(err, gc.ErrorMatches, `.*unknown version for series: "".*`)
}

func (s *supportedSeriesSuite) TestUbuntuSeriesVersion(c *gc.C) {
	isUbuntuTests := []struct {
		series   string
		expected string
	}{
		{"precise", "12.04"},
		{"raring", "13.04"},
		{"bionic", "18.04"},
		{"eoan", "19.10"},
		{"focal", "20.04"},
	}
	for _, v := range isUbuntuTests {
		ver, err := series.UbuntuSeriesVersion(v.series)
		c.Assert(err, gc.IsNil)
		c.Assert(ver, gc.Equals, v.expected)
	}
}

func (s *supportedSeriesSuite) TestUbuntuInvalidSeriesVersion(c *gc.C) {
	_, err := series.UbuntuSeriesVersion("firewolf")
	c.Assert(err, gc.ErrorMatches, `.*unknown version for series: "firewolf".*`)
}
