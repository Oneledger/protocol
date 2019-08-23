/*
	Copyright 2017-2018 OneLedger
*/
package version

import (
	"fmt"
)

type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	MetaData   string
}

var Fullnode *Version // Version of the source code
var Protocol *Version // Version of the protocol
var Client *Version   // Version of the client

// This should be the only copy of the version numbers, anywhere in the code.
func init() {
	// The protocol
	Protocol = NewVersion(0, 1, 2, "testnet", "Protocol")

	// The backend server (node) code
	Fullnode = NewVersion(0, 10, 9, "", "Fullnode")

	// Any of the clients used to connect
	Client = NewVersion(0, 10, 9, "", "Client")
}

func NewVersion(major, minor, patch int, release, meta string) *Version {
	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: release,
		MetaData:   meta,
	}
}

// Check for compatability, exact major, close to minor, bug fixes don't matter
func IsCompatible(base *Version, current *Version) bool {
	if base.Major != current.Major {
		return false
	}
	// Keep everyone within a couple of releases
	diff := base.Minor - current.Minor
	if diff < -2 || diff > 2 {
		return false
	}
	return true
}

func (v *Version) String() string {
	buffer := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		buffer += "-" + v.PreRelease
	}
	if v.MetaData != "" {
		buffer += "+" + v.MetaData
	}
	return buffer
}
