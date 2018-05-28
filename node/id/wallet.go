/*
	Copyright 2017 - 2018 OneLedger

	Group Mechanics
*/
package id

// A Group of users or groups.
type Wallet struct {
	Description string
	Accounts    []*Account
}
