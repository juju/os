// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package series_test

import (
	"github.com/juju/os"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/os/series"
)

type systemSuite struct {
	testing.CleanupSuite
}

var _ = gc.Suite(&systemSuite{})

func (s *systemSuite) TestSystemParsingToFromSeries(c *gc.C) {
	tests := []struct {
		system       series.System
		str          string
		parsedSystem series.System
	}{
		{series.System{OS: os.Ubuntu}, "system#os=ubuntu", series.System{OS: os.Ubuntu}},
		{series.System{OS: os.Ubuntu, Version: "focal"}, "focal", series.System{OS: os.Ubuntu, Version: "20.04"}},
		{series.System{OS: os.Ubuntu, Version: "18.04"}, "bionic", series.System{OS: os.Ubuntu, Version: "18.04"}},
		{series.System{OS: os.Windows}, "system#os=windows", series.System{OS: os.Windows}},
		{series.System{OS: os.Windows, Version: "win10"}, "win10", series.System{OS: os.Windows, Version: "win10"}},
		{series.System{OS: os.Ubuntu, Version: "18.04", Resource: "test"}, "system#os=ubuntu#version=18.04#resource=test", series.System{OS: os.Ubuntu, Version: "18.04", Resource: "test"}},
	}
	for i, v := range tests {
		str := v.system.String()
		c.Check(str, gc.Equals, v.str, gc.Commentf("test %d", i))
		s, err := series.ParseSystemFromSeries(str)
		c.Check(err, jc.ErrorIsNil)
		c.Check(s, jc.DeepEquals, v.parsedSystem)
	}
}
