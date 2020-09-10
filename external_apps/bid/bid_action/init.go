package bid_action

import (
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
	"github.com/Oneledger/protocol/serialize"
)

var BidAssetMap map[bid_data.BidAssetType]bid_data.BidAsset

const(
	//Bid
	BID_CREATE          action.Type = 0x101
	BID_CONTER_OFFER    action.Type = 0x102
	BID_CANCEL          action.Type = 0x103
	BID_BIDDER_DECISION action.Type = 0x104
	BID_EXPIRE          action.Type = 0x105
	BID_OWNER_DECISION  action.Type = 0x106
)

func init() {
	serialize.RegisterConcrete(new(CreateBid), "action_cb")
	serialize.RegisterConcrete(new(CancelBid), "action_cab")
	serialize.RegisterConcrete(new(BidderDecision), "action_bd")
	serialize.RegisterConcrete(new(OwnerDecision), "action_od")
	serialize.RegisterConcrete(new(CounterOffer), "action_co")
	serialize.RegisterConcrete(new(ExpireBid), "action_eb")
	action.RegisterTxType(BID_CREATE, "BID_CREATE")
	action.RegisterTxType(BID_CONTER_OFFER, "BID_CONTER_OFFER")
	action.RegisterTxType(BID_CANCEL, "BID_CANCEL")
	action.RegisterTxType(BID_BIDDER_DECISION, "BID_BIDDER_DECISION")
	action.RegisterTxType(BID_EXPIRE, "BID_EXPIRE")
	action.RegisterTxType(BID_OWNER_DECISION, "BID_OWNER_DECISION")
	BidAssetMap = make(map[bid_data.BidAssetType]bid_data.BidAsset)
	BidAssetMap[bid_data.BidAssetOns] = &bid_data.DomainAsset{}
	BidAssetMap[bid_data.BidAssetExample] = &bid_data.ExampleAsset{}
	fmt.Println("BidAssetMap in init: ", BidAssetMap)
}
