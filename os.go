// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

// Package os provides access to operating system related configuration.
package os

import (
	"github.com/juju/errors"
)

var HostOS = hostOS // for monkey patching

type OSType int

const (
	Unknown OSType = iota
	Ubuntu
	Windows
	OSX
	CentOS
	GenericLinux
	OpenSUSE
	Kubernetes
)

func (t OSType) String() string {
	switch t {
	case Ubuntu:
		return "Ubuntu"
	case Windows:
		return "Windows"
	case OSX:
		return "OSX"
	case CentOS:
		return "CentOS"
	case GenericLinux:
		return "GenericLinux"
	case OpenSUSE:
		return "OpenSUSE"
	case Kubernetes:
		return "Kubernetes"
	}
	return "Unknown"
}

// SystemString returns the OS type used by systems.
// NOTE: systems don't support kubernetes as an os type.
func (t OSType) SystemString() string {
	switch t {
	case Ubuntu:
		return "ubuntu"
	case Windows:
		return "windows"
	case OSX:
		return "osx"
	case CentOS:
		return "centos"
	case GenericLinux:
		return "genericlinux"
	case OpenSUSE:
		return "opensuse"
	}
	return ""
}

// EquivalentTo returns true if the OS type is equivalent to another
// OS type.
func (t OSType) EquivalentTo(t2 OSType) bool {
	if t == t2 {
		return true
	}
	return t.IsLinux() && t2.IsLinux()
}

// IsLinux returns true if the OS type is a Linux variant.
func (t OSType) IsLinux() bool {
	switch t {
	case Ubuntu, CentOS, GenericLinux, OpenSUSE:
		return true
	}
	return false
}

// ParseSystemOS parses OS type from a system to a OSType.
func ParseSystemOS(str string) (OSType, error) {
	switch str {
	case "ubuntu":
		return Ubuntu, nil
	case "windows":
		return Windows, nil
	case "osx":
		return OSX, nil
	case "centos":
		return CentOS, nil
	case "genericlinux":
		return GenericLinux, nil
	case "opensuse":
		return OpenSUSE, nil
	default:
		return Unknown, errors.NotValidf("system os type %q", str)
	}
}
