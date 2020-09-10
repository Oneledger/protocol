package bid_data

import (
	"github.com/Oneledger/protocol/action"
)

var _ BidAsset = &ExampleAsset{}

type ExampleAsset struct {
	ExampleName string `json:"exampleName"`
}

func (ea *ExampleAsset) ToString() string {
	return ea.ExampleName
}

func (ea *ExampleAsset) ValidateAsset(ctx *action.Context, owner action.Address) (bool, error) {
	return true, nil
}

func (ea *ExampleAsset) ExchangeAsset(ctx *action.Context, bidder action.Address, preOwner action.Address) (bool, error) {
	return true, nil
}

func (ea *ExampleAsset) NewAssetWithName(name string) BidAsset {
	asset := *ea
	asset.ExampleName = name
 	return &asset
}