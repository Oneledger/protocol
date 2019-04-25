/*

   ____             _          _
  / __ \           | |        | |
 | |  | |_ __   ___| | ___  __| | __ _  ___ _ __
 | |  | | '_ \ / _ \ |/ _ \/ _` |/ _` |/ _ \ '__|
 | |__| | | | |  __/ |  __/ (_| | (_| |  __/ |
  \____/|_| |_|\___|_|\___|\__,_|\__, |\___|_|
                                  __/ |
                                 |___/

	Copyright 2017 - 2019 OneLedger

*/

package balance

import (
	"errors"
	"math/big"

	"github.com/Oneledger/protocol/node/serialize"
)

type coinData struct {
	Currency Currency
	Amount   []byte
}

func (c *Coin) NewDataInstance() serialize.Data {
	return &coinData{}
}

func (c *Coin) Data() serialize.Data {
	return &coinData{c.Currency, c.Amount.Bytes()}
}

func (c *Coin) SetData(a interface{}) error {
	cd, ok := a.(*coinData)
	if !ok {
		return errors.New("Wrong data")
	}

	amt := new(big.Int)
	c.Currency = cd.Currency
	c.Amount = amt.SetBytes(cd.Amount)
	return nil
}

func (ad *coinData) SerialTag() string {
	return ""
}

func (cd *coinData) Primitive() serialize.DataAdapter {
	c := &Coin{}

	c.Currency = cd.Currency
	c.Amount = c.Amount.SetBytes(cd.Amount)

	return c
}
