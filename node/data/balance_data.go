package data

import (
	"errors"
	"github.com/Oneledger/protocol/node/serialize"
	"math/big"

	"github.com/Oneledger/protocol/node/log"
)

var (
	ErrWrongBalanceAdapter = errors.New("error in asserting to BalanceAdapter")
)

func init() {
	serialize.RegisterConcrete(new(BalanceData), TagBalanceData)
}

// BalanceData is an easy to serialize representation of a Balance object. A full Balance object can be recostructed
// from a BalanceAdapter object and vice versa.
// There is a map flattening of course for Coins
type BalanceData struct {
	Coins []CoinData `json:"pl"`
	Tag  string     `json:"tag"` // Tag is a field used to identify the type after ser/deser
	// will be useful in future
}

//

// CoinData is a flattening of coin map in a balance data type
type CoinData struct {
	CurId    int       `json:"curr_id"`
	CurName  string    `json:"curr_name"`
	CurChain ChainType `json:"curr_chain"`

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
	bd.Coins = make([]CoinData, 0, len(b.Amounts))


	for _, id := range b.coinOrder {
		coin := b.Amounts[id]
		cd := CoinData{
			CurId:    id,
			CurName:  coin.Currency.Name,
			CurChain: coin.Currency.Chain,
			Amount:   coin.Amount.Bytes(),
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

	b.Amounts = make(map[int]Coin)
	b.coinOrder = []int{}

	d := ba.Coins
	for i := range d {

		curID := d[i].CurId

		//convert string representation to big int
		amt := new(big.Int)
		amt = amt.SetBytes(d[i].Amount)

		coin := Coin{Amount: amt}
		coin.Currency.Id = curID
		coin.Currency.Name = d[i].CurName
		coin.Currency.Chain = d[i].CurChain

		b.Amounts[curID] = coin
		b.coinOrder = append(b.coinOrder, curID)
	}

	return nil
}

// Primitive gives the source object of the data through the adapter interface.
// Useful if you want to access the source object after deserialization.
func (bd *BalanceData) Primitive() serialize.DataAdapter {
	b := &Balance{}
	err := bd.extract(b)
	if err != nil {
		log.Error("error in get primitive of balance data", err)
		return nil
	}
	return b
}


func (bd *BalanceData) SerialTag() string {
	return bd.Tag
}