/*
	Copyright 2017 - 2018 OneLedger
*/
package data

// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}

type Coins []Coin
