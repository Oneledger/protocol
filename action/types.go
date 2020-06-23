package action

import (
	"encoding/json"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/tendermint/tendermint/abci/types"
)

// Address an action package over Address in data/keys package
type Address = keys.Address

// Balance an action package over Balance in data/balance
type Balance = balance.Balance

// Gas an action package over Gas in storage
type Gas = storage.Gas

// Amount is an easily serializable representation of coin. Nodes can create coin from the Amount object
// received over the network
type Amount struct {
	Currency string         `json:"currency"`
	Value    balance.Amount `json:"value"`
}

// New Amount creates a new amount account object
func NewAmount(currency string, amount balance.Amount) *Amount {
	return &Amount{currency, amount}
}

// IsValid checks the validity of the currency and the amount string in the account object, which may be received
// over a network.
func (a Amount) IsValid(list *balance.CurrencySet) bool {
	currency, ok := list.GetCurrencyByName(a.Currency)
	if !ok {
		return false
	}

	coin := currency.NewCoinFromAmount(a.Value)
	return coin.IsValid()
}

// String returns a string representation of the Amount object.
func (a Amount) String() string {
	result, _ := json.Marshal(a)
	return string(result)
}

// ToCoin converts an easier to transport Amount object to a Coin object in Oneledger protocol.
// It takes the action context to determine the currency from which to create the coin.
func (a Amount) ToCoin(list *balance.CurrencySet) balance.Coin {

	// get currency of Amount a
	currency, ok := list.GetCurrencyByName(a.Currency)
	if !ok {
		return balance.Coin{}
	}

	// parse float string
	return currency.NewCoinFromAmount(a.Value)
}

func (a Amount) ToCoinWithBase(list *balance.CurrencySet) balance.Coin {

	// get currency of Amount a
	currency, ok := list.GetCurrencyByName(a.Currency)
	if !ok {
		return balance.Coin{}
	}

	// parse float string
	return currency.NewCoinFromInt(a.Value.BigInt().Int64())
}

type Response struct {
	Data      []byte
	Events    []types.Event
	Log       string
	Info      string
	GasWanted int64
	GasUsed   int64
}
