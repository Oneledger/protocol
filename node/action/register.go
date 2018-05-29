/*
	Copyright 2017 - 2018 OneLedger

	Register this identity with the other nodes. As an externl identity
*/
package action

// Synchronize a swap between two users
type Register struct {
	Base

	Identity string
	Node     string
}
