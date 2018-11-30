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

var Current *Version  // Version of the source code
var Protocol *Version // Version of the protocol

// This should be the only copy of the version numbers, anywhere in the code.
func init() {
	Current = &Version{
		Major:      0,
		Minor:      6,
		Patch:      2,
		PreRelease: "",
		MetaData:   "",
	}

	// The code vs. the underlying protocol version. They will drift at some point...
	Protocol = Current
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
