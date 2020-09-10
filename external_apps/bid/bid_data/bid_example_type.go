package bid_data

import (
	"github.com/Oneledger/protocol/action"
	"strconv"
)

var _ BidAsset = &ExampleAsset{}

type ExampleAsset struct {
	ExampleField int `json:"exampleField"`
}

func NewExampleAsset(exampleValue int) *ExampleAsset {
	return &ExampleAsset{
		ExampleField: exampleValue,
	}
}

func (ta *ExampleAsset) ToString() string {
	return strconv.Itoa(ta.ExampleField)
}

func (ta *ExampleAsset) ValidateAsset(ctx *action.Context, owner action.Address) (bool, error) {
	return true, nil
}

func (ta *ExampleAsset) ExchangeAsset(ctx *action.Context, bidder action.Address, preOwner action.Address) (bool, error) {
	return true, nil
}

func (ta *ExampleAsset) SetName(name string) {

}