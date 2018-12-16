/*
	Copyright 2017-2018 OneLedger
*/
package version

import (
	"fmt"

	"github.com/Oneledger/protocol/node/serial"
)

type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	MetaData   string
}

func init() {
	serial.Register(Version{})
}

var Fullnode *Version // Version of the source code
var Protocol *Version // Version of the protocol
var Client *Version   // Version of the protocol

// This should be the only copy of the version numbers, anywhere in the code.
func init() {
	// The protocol
	Protocol = NewVersion(0, 1, 0, "testnet", "Protocol")

	// The backend server (node) code
	Fullnode = NewVersion(0, 8, 0, "", "Fullnode")

	// Any of the clients used to connect
	Client = NewVersion(0, 8, 0, "", "Client")
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
