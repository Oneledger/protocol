package bid_action

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

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
}

//todo old stuff, will need to migrate to new structure
//func EnableBidding(r action.Router) error {
//	err := r.AddHandler(action.BID_CREATE, CreateBid{})
//	if err != nil {
//		return errors.Wrap(err, "CreateBidTx")
//	}
//	err = r.AddHandler(action.BID_CANCEL, CancelBid{})
//	if err != nil {
//		return errors.Wrap(err, "cancelBidTx")
//	}
//	err = r.AddHandler(action.BID_EXPIRE, ExpireBid{})
//	if err != nil {
//		return errors.Wrap(err, "expireBidTx")
//	}
//	err = r.AddHandler(action.BID_BIDDER_DECISION, BidderDecision{})
//	if err != nil {
//		return errors.Wrap(err, "bidderDecisionTx")
//	}
//	err = r.AddHandler(action.BID_CONTER_OFFER, CounterOffer{})
//	if err != nil {
//		return errors.Wrap(err, "counterOfferTx")
//	}
//	err = r.AddHandler(action.BID_OWNER_DECISION, OwnerDecision{})
//	if err != nil {
//		return errors.Wrap(err, "ownerDecisionTx")
//	}
//
//	return nil
//}
//
//func EnableInternalBidding(r action.Router) error {
//	err := r.AddHandler(action.BID_EXPIRE, ExpireBid{})
//	if err != nil {
//		return errors.Wrap(err, "ExpireBidConvsTx")
//	}
//	return nil
//}