// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series_test

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/os"
	"github.com/juju/os/series"
)

func (s *supportedSeriesSuite) TestSeriesVersion(c *gc.C) {
	// There is no distro-info on Windows or CentOS.
	if os.HostOS() != os.Ubuntu {
		c.Skip("This test is only relevant on Ubuntu.")
	}
	vers, err := series.SeriesVersion("precise")
	if err != nil && err.Error() == `invalid series "precise"` {
		c.Fatalf(`Unable to lookup series "precise", you may need to: apt-get install distro-info`)
	}
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(vers, gc.Equals, "12.04")
}

func (s *supportedSeriesSuite) TestSupportedSeries(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.UbuntuDistroInfoPath, filename)

	expectedSeries := []string{"artful", "bionic", "cosmic", "disco", "eoan", "focal", "precise", "quantal", "raring", "saucy", "trusty", "utopic", "vivid", "wily", "xenial", "yakkety", "zesty"}
	series := series.SupportedSeries()
	sort.Strings(series)
	c.Assert(series, gc.DeepEquals, expectedSeries)
}

func (s *supportedSeriesSuite) TestUpdateSeriesVersions(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.UbuntuDistroInfoPath, filename)

	expectedSeries := []string{"artful", "bionic", "cosmic", "disco", "eoan", "focal", "precise", "quantal", "raring", "saucy", "trusty", "utopic", "vivid", "wily", "xenial", "yakkety", "zesty"}
	checkSeries := func() {
		series := series.SupportedSeries()
		sort.Strings(series)
		c.Assert(series, gc.DeepEquals, expectedSeries)
	}
	checkSeries()

	// Updating the file does not normally trigger an update;
	// we only refresh automatically one time. After that, we
	// must explicitly refresh.
	err = ioutil.WriteFile(filename, []byte(distInfoData2), 0644)
	c.Assert(err, jc.ErrorIsNil)
	checkSeries()

	expectedSeries = append([]string{"ornery"}, expectedSeries...)
	sort.Strings(expectedSeries)
	series.UpdateSeriesVersions()
	checkSeries()
}

func (s *supportedSeriesSuite) TestESMSupportedJujuSeries(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.UbuntuDistroInfoPath, filename)

	expectedSeries := []string{"focal", "bionic", "xenial", "trusty"}
	series := series.ESMSupportedJujuSeries()
	c.Assert(series, jc.DeepEquals, expectedSeries)
}

func (s *supportedSeriesSuite) TestOSSeries(c *gc.C) {
	restore := series.HideUbuntuSeries()
	defer restore()

	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.UbuntuDistroInfoPath, filename)

	osType, err := series.GetOSFromSeries("raring")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(osType, gc.Equals, os.Ubuntu)
}

func (s *supportedSeriesSuite) TestSupportedJujuControllerSeries(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.UbuntuDistroInfoPath, filename)

	expectedSeries := []string{"focal", "bionic", "xenial"}
	series := series.SupportedJujuControllerSeries()
	c.Assert(series, jc.DeepEquals, expectedSeries)
}

func (s *supportedSeriesSuite) TestSupportedJujuWorkloadSeries(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.UbuntuDistroInfoPath, filename)

	expectedSeries := []string{"focal", "bionic", "xenial", "centos7", "centos8", "genericlinux", "kubernetes", "opensuseleap", "win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2", "win2016", "win2016hv", "win2016nano", "win2019", "win7", "win8", "win81"}
	series := series.SupportedJujuWorkloadSeries()
	c.Assert(series, jc.DeepEquals, expectedSeries)
}

func (s *supportedSeriesSuite) TestSupportedJujuSeries(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "ubuntu.csv")
	err := ioutil.WriteFile(filename, []byte(distInfoData), 0644)
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(series.UbuntuDistroInfoPath, filename)

	expectedSeries := []string{"focal", "bionic", "xenial", "centos7", "centos8", "genericlinux", "kubernetes", "opensuseleap", "win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2", "win2016", "win2016hv", "win2016nano", "win2019", "win7", "win8", "win81"}
	series := series.SupportedJujuSeries()
	c.Assert(series, jc.DeepEquals, expectedSeries)
}

func (s *supportedSeriesSuite) TestLatestLts(c *gc.C) {
	table := []struct {
		latest, want string
	}{
		{"testseries", "testseries"},
		{"", "focal"},
	}
	for _, test := range table {
		series.SetLatestLtsForTesting(test.latest)
		got := series.LatestLts()
		c.Assert(got, gc.Equals, test.want)
	}
}

func (s *supportedSeriesSuite) TestSetLatestLtsForTesting(c *gc.C) {
	table := []struct {
		value, want string
	}{
		{"1", "focal"}, {"2", "1"}, {"3", "2"}, {"4", "3"},
	}
	for _, test := range table {
		got := series.SetLatestLtsForTesting(test.value)
		c.Assert(got, gc.Equals, test.want)
	}
}

func (s *supportedSeriesSuite) TestSupportedLts(c *gc.C) {
	got := series.SupportedLts()
	want := []string{"xenial", "bionic", "focal"}
	c.Assert(got, gc.DeepEquals, want)
}

const distInfoData = `version,codename,series,created,release,eol,eol-server,eol-esm
4.10,Warty Warthog,warty,2004-03-05,2004-10-20,2006-04-30
5.04,Hoary Hedgehog,hoary,2004-10-20,2005-04-08,2006-10-31
5.10,Breezy Badger,breezy,2005-04-08,2005-10-12,2007-04-13
6.06 LTS,Dapper Drake,dapper,2005-10-12,2006-06-01,2009-07-14,2011-06-01
6.10,Edgy Eft,edgy,2006-06-01,2006-10-26,2008-04-25
7.04,Feisty Fawn,feisty,2006-10-26,2007-04-19,2008-10-19
7.10,Gutsy Gibbon,gutsy,2007-04-19,2007-10-18,2009-04-18
8.04 LTS,Hardy Heron,hardy,2007-10-18,2008-04-24,2011-05-12,2013-05-09
8.10,Intrepid Ibex,intrepid,2008-04-24,2008-10-30,2010-04-30
9.04,Jaunty Jackalope,jaunty,2008-10-30,2009-04-23,2010-10-23
9.10,Karmic Koala,karmic,2009-04-23,2009-10-29,2011-04-29
10.04 LTS,Lucid Lynx,lucid,2009-10-29,2010-04-29,2013-05-09,2015-04-29
10.10,Maverick Meerkat,maverick,2010-04-29,2010-10-10,2012-04-10
11.04,Natty Narwhal,natty,2010-10-10,2011-04-28,2012-10-28
11.10,Oneiric Ocelot,oneiric,2011-04-28,2011-10-13,2013-05-09
12.04 LTS,Precise Pangolin,precise,2011-10-13,2012-04-26,2017-04-26,2017-04-26,2019-04-26
12.10,Quantal Quetzal,quantal,2012-04-26,2012-10-18,2014-05-16
13.04,Raring Ringtail,raring,2012-10-18,2013-04-25,2014-01-27
13.10,Saucy Salamander,saucy,2013-04-25,2013-10-17,2014-07-17
14.04 LTS,Trusty Tahr,trusty,2013-10-17,2014-04-17,2019-04-17,2019-04-17,2022-04-17
14.10,Utopic Unicorn,utopic,2014-04-17,2014-10-23,2015-07-23
15.04,Vivid Vervet,vivid,2014-10-23,2015-04-23,2016-01-23
15.10,Wily Werewolf,wily,2015-04-23,2015-10-22,2016-07-22
16.04 LTS,Xenial Xerus,xenial,2015-10-22,2016-04-21,2021-04-21,2021-04-21,2024-04-21
16.10,Yakkety Yak,yakkety,2016-04-21,2016-10-13,2017-07-20
17.04,Zesty Zapus,zesty,2016-10-13,2017-04-13,2018-01-13
17.10,Artful Aardvark,artful,2017-04-13,2017-10-19,2018-07-19
18.04 LTS,Bionic Beaver,bionic,2017-10-19,2018-04-26,2023-04-26,2023-04-26,2028-04-26
18.10,Cosmic Cuttlefish,cosmic,2018-04-26,2018-10-18,2019-07-18
19.04,Disco Dingo,disco,2018-10-18,2019-04-18,2020-01-18
19.10,Eoan Ermine,eoan,2019-04-18,2019-10-17,2020-07-17
20.04 LTS,Focal Fossa,focal,2019-10-17,2020-04-23,2025-04-23,2025-04-23,2030-04-23
`

const distInfoData2 = distInfoData + `
14.04 LTS,Firewolf,firewolf,2013-10-17,2014-04-17
94.04 LTS,Ornery Omega,ornery,2094-10-17,2094-04-17,2099-04-17
`

type isolationSupportedSeriesSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&isolationSupportedSeriesSuite{})

func (s *isolationSupportedSeriesSuite) TestBadFilePath(c *gc.C) {
	d := c.MkDir()
	filename := filepath.Join(d, "bad-file.csv")
	s.PatchValue(series.UbuntuDistroInfoPath, filename)

	expectedSeries := []string{"artful", "bionic", "centos7", "centos8", "cosmic", "disco", "eoan", "focal", "genericlinux", "groovy", "opensuseleap", "precise", "quantal", "raring", "saucy", "trusty", "utopic", "vivid", "wily", "win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2", "win2016", "win2016hv", "win2016nano", "win2019", "win7", "win8", "win81", "xenial", "yakkety", "zesty"}
	series := series.SupportedSeries()
	sort.Strings(series)
	c.Assert(series, gc.DeepEquals, expectedSeries)
}
