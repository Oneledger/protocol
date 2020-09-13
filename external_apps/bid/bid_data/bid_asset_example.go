package bid_data

import (
	"github.com/Oneledger/protocol/action"
)

var _ BidAsset = &ExampleAsset{}

type ExampleAsset struct {
	exampleName string
}

func (ea *ExampleAsset) ToString() string {
	return ea.exampleName
}

func (ea *ExampleAsset) ValidateAsset(ctx *action.Context, owner action.Address) (bool, error) {
	return true, nil
}

func (ea *ExampleAsset) ExchangeAsset(ctx *action.Context, bidder action.Address, preOwner action.Address) (bool, error) {
	return true, nil
}

func (ea *ExampleAsset) NewAssetWithName(name string) BidAsset {
	asset := *ea
	asset.exampleName = name
	return &asset
}
