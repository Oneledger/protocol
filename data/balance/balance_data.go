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
	"sort"

	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/serialize"
)

// BalanceData is an easy to serialize representation of a Balance object. A full Balance object can be recostructed
// from a BalanceAdapter object and vice versa.
// There is a map flattening of course for Coins
type BalanceData struct {
	Coins []coinData `json:"coins"`
	Tag   string     `json:"tag"` // Tag is a field used to identify the type after ser/deser
	// will be useful in future
}

// coinData is a flattening of coin map in a balance data type
type coinData struct {
	CurName    string     `json:"curr_name"`
	CurChain   chain.Type `json:"curr_chain"`
	CurDecimal int64      `json:"curr_decimal"`

	Amount []byte `json:"amt"`
}

//
func init() {}

func (b *Balance) NewDataInstance() serialize.Data {
	return &BalanceData{}
}

// Data creates a BalanceData from a given Balance object,
// the coins are flattened to a list in the generator itself
// ideally there should be no change done to a data after this step. This datatype can go straight to serialization.
func (b *Balance) Data() serialize.Data {
	//initialize with source pointer
	bd := &BalanceData{Tag: "balance_data"}
	// this allows to reserve capacity so the process of adding
	// items to the list
	bd.Coins = make([]coinData, 0, len(b.Amounts))

	currencyList := []string{}
	for key := range b.Amounts {
		currencyList = append(currencyList, key)
	}

	sort.Strings(currencyList)

	for _, key := range currencyList {
		coin := b.Amounts[key]
		cd := coinData{
			CurName:    coin.Currency.Name,
			CurChain:   coin.Currency.Chain,
			CurDecimal: coin.Currency.Decimal,
			Amount:     coin.Amount.Bytes(),
		}

		bd.Coins = append(bd.Coins, cd)
	}
	return bd
}

// SetData sets the balance object back from a BalanceData object
func (b *Balance) SetData(obj interface{}) error {
	ba, ok := obj.(*BalanceData)
	if !ok {
		return ErrWrongBalanceAdapter
	}
	return ba.extract(b)
}

//

// Extract recreates the Balance object form the info BalanceData holds after deserialization/
func (ba *BalanceData) extract(b *Balance) error {

	b.Amounts = make(map[string]Coin)

	d := ba.Coins
	for i := range d {

		//convert string representation to big int
		amt := NewAmount(0)
		amt = amt.SetBytes(d[i].Amount)

		coin := Coin{Amount: amt}
		coin.Currency.Name = d[i].CurName
		coin.Currency.Chain = d[i].CurChain
		coin.Currency.Decimal = d[i].CurDecimal

		b.Amounts[coin.Currency.StringKey()] = coin
	}

	return nil
}

func (bd *BalanceData) SerialTag() string {
	return bd.Tag
}
