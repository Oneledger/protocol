/*
	Copyright 2017-2018 OneLedger
*/
package version

import "fmt"

type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	MetaData   string
}

var Current *Version

// This should be the only copy of the version numbers, anywhere in the code.
func init() {
	Current = &Version{
		Major:      0,
		Minor:      5,
		Patch:      0,
		PreRelease: "",
		MetaData:   "",
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
