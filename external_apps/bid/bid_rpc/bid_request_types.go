package bid_rpc

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
)

type ListBidConvRequest struct {
	BidConvId bid_data.BidConvId `json:"bidConvId"`
}

type BidConvStat struct {
	BidConv bid_data.BidConv    `json:"bidConv"`
	ActiveOffer bid_data.BidOffer `json:"activeOffer"`
	InactiveOffers  []bid_data.BidOffer `json:"inactiveOffers"`
}

type ListBidConvsReply struct {
	BidConvStats []BidConvStat `json:"bidConvStats"`
	Height       int64         `json:"height"`
}

type ListBidConvsRequest struct {
	State     bid_data.BidConvState `json:"state"`
	Owner     keys.Address          `json:"owner"`
	AssetName string                `json:"assetName"`
	AssetType bid_data.BidAssetType `json:"assetType"`
	Bidder    keys.Address          `json:"bidder"`
}

type CreateBidRequest struct {
	BidConvId  bid_data.BidConvId    `json:"bidConvId"`
	AssetOwner keys.Address          `json:"assetOwner"`
	AssetName  string                `json:"assetName"`
	AssetType  bid_data.BidAssetType `json:"assetType"`
	Bidder     keys.Address          `json:"bidder"`
	Amount     action.Amount         `json:"amount"`
	Deadline   int64                 `json:"deadline"`
	GasPrice   action.Amount         `json:"gasPrice"`
	Gas        int64                 `json:"gas"`
}

type CounterOfferRequest struct {
	BidConvId  bid_data.BidConvId `json:"bidConvId"`
	AssetOwner keys.Address       `json:"assetOwner"`
	Amount     action.Amount      `json:"amount"`
	GasPrice   action.Amount      `json:"gasPrice"`
	Gas        int64              `json:"gas"`
}

type CancelBidRequest struct {
	BidConvId bid_data.BidConvId `json:"bidConvId"`
	Bidder    keys.Address       `json:"bidder"`
	GasPrice  action.Amount      `json:"gasPrice"`
	Gas       int64              `json:"gas"`
}

type OwnerDecisionRequest struct {
	BidConvId bid_data.BidConvId   `json:"bidConvId"`
	Owner     keys.Address         `json:"owner"`
	Decision  bid_data.BidDecision `json:"decision"`
	GasPrice  action.Amount        `json:"gasPrice"`
	Gas       int64                `json:"gas"`
}

type BidderDecisionRequest struct {
	BidConvId bid_data.BidConvId   `json:"bidConvId"`
	Bidder    keys.Address         `json:"bidder"`
	Decision  bid_data.BidDecision `json:"decision"`
	GasPrice  action.Amount        `json:"gasPrice"`
	Gas       int64                `json:"gas"`
}
