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

	"github.com/Oneledger/protocol/serialize"
)

type CoinData struct {
	Currency Currency
	Amount   []byte
}

func (c *Coin) NewDataInstance() serialize.Data {
	return &CoinData{}
}

func (c *Coin) Data() serialize.Data {
	return &CoinData{c.Currency, c.Amount.Bytes()}
}

func (c *Coin) SetData(a interface{}) error {
	cd, ok := a.(*CoinData)
	if !ok {
		return errors.New("Wrong data")
	}

	amt := new(big.Int)
	c.Currency = cd.Currency
	c.Amount = amt.SetBytes(cd.Amount)
	return nil
}

func (ad *CoinData) SerialTag() string {
	return ""
}
