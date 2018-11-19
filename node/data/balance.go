/*
	Copyright 2017-2018 OneLedger
*/

package data

import (
	"fmt"

	"github.com/Oneledger/protocol/node/serial"
)

// Wrap the amount with owner information
type Balance struct {
	// Address id.Address
	Amount Coin
}

func init() {
	serial.Register(Balance{})
}

// TODO: Should return a pointer, as per Go conventions
func NewBalance(amount int64, currency string) Balance {
	return Balance{Amount: NewCoin(amount, currency)}
}

//String used in fmt and Dump
func (balance Balance) String() string {
	return fmt.Sprintf("%s %s", balance.Amount, balance.Amount.Currency.Name)
}
