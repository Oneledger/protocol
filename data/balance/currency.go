/*

   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/

	Copyright 2017 - 2019 OneLedger

*/

package balance

import (
	"math/big"

	"github.com/Oneledger/protocol/data/chain"
)

/*
 Currency starts here
*/

type Currency struct {
	Name  string     `json:"name"`
	Chain chain.Type `json:"chain"`

	Decimal int64 `json:"decimal"`
}

func (c Currency) Base() *big.Int {
	return big.NewInt(0).Exp(big.NewInt(10), big.NewInt(c.Decimal), nil)
}

// Create a coin from integer (not fractional)
func (c Currency) NewCoinFromInt(amount int64) Coin {
	return Coin{
		Currency: c,
		Amount:   big.NewInt(amount),
	}
}
