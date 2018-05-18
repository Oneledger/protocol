/*
	Copyright 2017 - 2018 OneLedger

	Basic datatypes
*/
package data

import ()

type Chain struct {
}

type ChainNode struct {
	// TODO: How to navigate to the node via grpc
}

type Balance struct {
	// Address id.Address
	Amount Coin
}
