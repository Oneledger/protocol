package query

import (
	"errors"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/bidding"
	codes "github.com/Oneledger/protocol/status_codes"
)

// show single bid conversation by id
func (svc *Service) ShowBidConv(req client.ListBidConvRequest, reply *client.ListBidConvsReply) error {
	bidConv, _, err := svc.bidMaster.BidConv.QueryAllStores(req.BidConvId)
	if err != nil {
		svc.logger.Error("error getting bid conversation", err)
		return codes.ErrGettingBidConv
	}

	bidOffers := svc.bidMaster.BidOffer.GetOffers(bidConv.BidConvId, bidding.BidOfferInvalid, bidding.TypeInvalid)

	bcs := client.BidConvStat{
		BidConv: *bidConv,
		Offers:  bidOffers,
	}

	*reply = client.ListBidConvsReply{
		BidConvStats: []client.BidConvStat{bcs},
		Height:       svc.proposalMaster.Proposal.GetState().Version(),
	}
	return nil
}

// list single proposal by id or list proposals
func (svc *Service) ListBidConvs(req client.ListBidConvsRequest, reply *client.ListBidConvsReply) error {
	// Validate parameters
	if len(req.Owner) != 0 {
		err := req.Owner.Err()
		if err != nil {
			return errors.New("invalid asset owner address")
		}
	}

	if len(req.Bidder) != 0 {
		err := req.Bidder.Err()
		if err != nil {
			return errors.New("invalid asset bidder address")
		}
	}

	// Query in single store if specified
	bms := svc.bidMaster
	var bidConvs []bidding.BidConv
	if req.State != bidding.BidStateInvalid {
		bidConvs = bms.BidConv.FilterBidConvs(req.State, req.Owner, req.AssetName, req.AssetType, req.Bidder)
	} else { // Query in all stores otherwise
		active := bms.BidConv.FilterBidConvs(bidding.BidStateActive, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		succeed := bms.BidConv.FilterBidConvs(bidding.BidStateSucceed, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		rejected := bms.BidConv.FilterBidConvs(bidding.BidStateRejected, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		expired := bms.BidConv.FilterBidConvs(bidding.BidStateExpired, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		cancelled := bms.BidConv.FilterBidConvs(bidding.BidStateCancelled, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		bidConvs = append(bidConvs, active...)
		bidConvs = append(bidConvs, succeed...)
		bidConvs = append(bidConvs, rejected...)
		bidConvs = append(bidConvs, expired...)
		bidConvs = append(bidConvs, cancelled...)
	}

	// Organize reply packet:
	// Bid conversations and their offers
	bidConvStats := make([]client.BidConvStat, len(bidConvs))
	for i, bidConv := range bidConvs {
		bidOffers := svc.bidMaster.BidOffer.GetOffers(bidConv.BidConvId, bidding.BidOfferInvalid, bidding.TypeInvalid)
		bcs := client.BidConvStat{
			BidConv: bidConv,
			Offers: bidOffers,
		}
		bidConvStats[i] = bcs
	}

	*reply = client.ListBidConvsReply{
		BidConvStats: bidConvStats,
		Height:       svc.proposalMaster.Proposal.GetState().Version(),
	}
	return nil
}

// list active offers
func (svc *Service) ListActiveOffers(req client.ListActiveOffersRequest, reply *client.ListActiveOffersReply) error {
	// Validate parameters
	if len(req.Owner) != 0 {
		err := req.Owner.Err()
		if err != nil {
			return errors.New("invalid asset owner address")
		}
	}

	if len(req.Bidder) != 0 {
		err := req.Bidder.Err()
		if err != nil {
			return errors.New("invalid asset bidder address")
		}
	}

	// get all active offers
	bms := svc.bidMaster
	offers := bms.BidOffer.GetOffers("", bidding.BidOfferActive, req.OfferType)
	activeOfferStats := make([]client.ActiveOfferStat, len(offers))
	for i, offer := range offers {
		// get corresponding bid conversation to show the detail
		bidConv, err := svc.bidMaster.BidConv.WithPrefixType(bidding.BidStateActive).Get(offer.BidConvId)
		if err != nil {
			svc.logger.Error("error getting bid conversation", err)
			return codes.ErrGettingBidConv
		}
		if len(req.Bidder) != 0 && !req.Bidder.Equal(bidConv.Bidder) {
			continue
		}
		if len(req.Owner) != 0 && !req.Owner.Equal(bidConv.AssetOwner) {
			continue
		}
		if req.AssetType != bidding.BidAssetInvalid && req.AssetType != bidConv.AssetType {
			continue
		}
		if req.AssetName != bidConv.Asset.ToString() {
			continue
		}
		aos := client.ActiveOfferStat{
			ActiveOffer: offer,
			BidConv: *bidConv,
		}
		activeOfferStats[i] = aos
	}

	*reply = client.ListActiveOffersReply{
		ActiveOffers: activeOfferStats,
		Height:       svc.proposalMaster.Proposal.GetState().Version(),
	}
	return nil
}


