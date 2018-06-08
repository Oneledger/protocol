/*
	Copyright 2017 - 2018 OneLedger

	Basic datatypes
*/
package data

/*
type Chain struct {
}

type ChainNode struct {
	// TODO: How to navigate to the node via grpc
}
*/

type Balance struct {
	// Address id.Address
	Amount Coin
}

func NewBalance(amount int64, currency string) Balance {
	return Balance{Amount: NewCoin(amount, currency)}
}
