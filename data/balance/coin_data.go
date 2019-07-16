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

	"github.com/Oneledger/protocol/serialize"
)

type CoinData struct {
	Currency Currency `json:"currency"`
	Amount   []byte   `json:"amount"`
}

func (c *Coin) NewDataInstance() serialize.Data {
	return &CoinData{}
}

func (c *Coin) Data() serialize.Data {
	b, _ := c.Amount.MarshalJSON()
	return &CoinData{c.Currency, b}
}

func (c *Coin) SetData(a interface{}) error {
	cd, ok := a.(*CoinData)
	if !ok {
		return errors.New("Wrong coin data")
	}

	amt := &Amount{}
	err := amt.UnmarshalJSON(cd.Amount)
	if err != nil {
		return err
	}
	c.Currency = cd.Currency
	c.Amount = amt
	return nil
}

func (ad *CoinData) SerialTag() string {
	return ""
}
