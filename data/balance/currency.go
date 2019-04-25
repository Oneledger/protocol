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
	"encoding/hex"
	"math/big"

	"golang.org/x/crypto/ripemd160"

	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/serialize"
)

/*
 Currency starts here
*/


func GetCurrencies() map[string]Currency {
	return currencies
}

type Currency struct {
	Name  string          `json:"name"`
	Chain chain.ChainType `json:"chain"`
	Id    int             `json:"id"`
}

// Look up the currency
func NewCurrency(currency string) Currency {
	return currencies[currency]
}

func GetBase(currency string) *big.Float {
	return GetExtra(currency).Units
}

// Key sets a encodable key for the currency entry, we may end up using currencyCodes instead.
func (c Currency) Key() (string, error) {
	hasher := ripemd160.New()

	buffer, err := serialize.JSONSzr.Serialize(c)
	if err != nil {
		log.Fatal("hash serialize failed", "err", err)
	}

	_, err = hasher.Write(buffer)
	if err != nil {
		log.Fatal("hasher failed", "err", err)
	}

	buffer = hasher.Sum(nil)

	return hex.EncodeToString(buffer), nil
}

/*
	Currency Extra
*/
type Extra struct {
	Units   *big.Float
	Decimal int
	Format  uint8
}

// TODO: Separated from Currency to avoid serializing big floats and giving out this info


func GetExtra(currency string) Extra {
	if value, ok := currenciesExtra[currency]; ok {
		return value
	}
	return currenciesExtra["OLT"]
}
