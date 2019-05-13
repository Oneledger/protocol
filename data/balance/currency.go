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
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
	"math"
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

func (c Currency) StringKey() string {
	return hex.EncodeToString(c.Bytes())
}

func (c Currency) Bytes() []byte {

	dat, _ := json.Marshal(c)
	r := ripemd160.New()
	r.Write(dat)
	return r.Sum(nil)
}

// TODO
// Create a coin from float
func (c Currency) NewCoinFromFloat64(amount float64) Coin {

	base := math.Pow10(int(c.Decimal))

	amountBigFloat := new(big.Float)
	amountBigFloat.SetFloat64(amount)
	// Set precision if required.
	// amountBigFloat.SetPrec(64)

	baseFloat := new(big.Float)
	baseFloat.SetFloat64(base)

	amountBigFloat.Mul(amountBigFloat, baseFloat)

	result := new(big.Int)
	amountBigFloat.Int(result) // store converted number in result

	return Coin{
		Currency: c,
		Amount:   result,
	}
}

// Create a coin from bytes, the bytes must come from Big.Int.
func (c Currency) NewCoinFromBytes(amount []byte) Coin {
	return Coin{
		Currency: c,
		Amount:   big.NewInt(0).SetBytes(amount),
	}
}
