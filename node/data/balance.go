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

func NewBalance(amount int64, currency string) Balance {
	return Balance{Amount: NewCoin(amount, currency)}
}

func (balance Balance) AsString() string {
	return fmt.Sprintf("%s %s", balance.Amount.AsString(), balance.Amount.Currency.Name)
}
