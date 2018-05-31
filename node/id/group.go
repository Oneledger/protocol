/*
	Copyright 2017 - 2018 OneLedger

	Group Mechanics
*/
package id

// The ability to execute a specific peice of code
type RolePrivledge struct {
	Name string
}

// The ability to use a specific peice of data
type DataPrivledge struct {
	Name string
}

// A Group of users or groups.
type Group struct {
	Name  string
	Roles []RolePrivledge
	Data  []DataPrivledge

	Identities []*Identity
	Groups     []*Group
}
