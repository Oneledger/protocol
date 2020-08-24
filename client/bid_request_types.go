package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/bidding"
	"github.com/Oneledger/protocol/data/keys"
)

type ListBidConvRequest struct {
	BidConvId bidding.BidConvId `json:"bidConvId"`
}

type BidConvStat struct {
	BidConv bidding.BidConv    `json:"bidConv"`
	Offers  []bidding.BidOffer `json:"bidOffers"`
}

type ActiveOfferStat struct {
	ActiveOffer bidding.BidOffer `json:"activeOffer"`
	BidConv     bidding.BidConv  `json:"bidConv"`
}

type ListBidConvsReply struct {
	BidConvStats []BidConvStat `json:"bidConvStats"`
	Height       int64         `json:"height"`
}

type ListBidConvsRequest struct {
	State     bidding.BidConvState `json:"state"`
	Owner     keys.Address         `json:"owner"`
	AssetName string               `json:"assetName"`
	AssetType bidding.BidAssetType `json:"assetType"`
	Bidder    keys.Address         `json:"bidder"`
}

type ListActiveOffersRequest struct {
	Owner     keys.Address         `json:"owner"`
	AssetName string               `json:"assetName"`
	AssetType bidding.BidAssetType `json:"assetType"`
	Bidder    keys.Address         `json:"bidder"`
	OfferType bidding.BidOfferType `json:"offerType"`
}

type ListActiveOffersReply struct {
	ActiveOffers []ActiveOfferStat `json:"activeOffers"`
	Height       int64             `json:"height"`
}

type CreateBidRequest struct {
	BidConvId  bidding.BidConvId    `json:"bidConvId"`
	AssetOwner keys.Address         `json:"assetOwner"`
	Asset      bidding.BidAsset     `json:"asset"`
	AssetType  bidding.BidAssetType `json:"assetType"`
	Bidder     keys.Address         `json:"bidder"`
	Amount     action.Amount        `json:"amount"`
	Deadline   int64                `json:"deadline"`
	GasPrice   action.Amount        `json:"gasPrice"`
	Gas        int64                `json:"gas"`
}

type CounterOfferRequest struct {
	BidConvId  bidding.BidConvId `json:"bidConvId"`
	AssetOwner keys.Address      `json:"assetOwner"`
	Amount     action.Amount     `json:"amount"`
	GasPrice   action.Amount     `json:"gasPrice"`
	Gas        int64             `json:"gas"`
}

type CancelBidRequest struct {
	BidConvId bidding.BidConvId `json:"bidConvId"`
	Bidder    keys.Address      `json:"bidder"`
	GasPrice  action.Amount     `json:"gasPrice"`
	Gas       int64             `json:"gas"`
}

type OwnerDecisionRequest struct {
	BidConvId bidding.BidConvId   `json:"bidConvId"`
	Owner     keys.Address        `json:"owner"`
	Decision  bidding.BidDecision `json:"decision"`
	GasPrice  action.Amount       `json:"gasPrice"`
	Gas       int64               `json:"gas"`
}

type BidderDecisionRequest struct {
	BidConvId bidding.BidConvId   `json:"bidConvId"`
	Bidder    keys.Address        `json:"bidder"`
	Decision  bidding.BidDecision `json:"decision"`
	GasPrice  action.Amount       `json:"gasPrice"`
	Gas       int64               `json:"gas"`
}
