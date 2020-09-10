package bid_rpc_query

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
	"github.com/Oneledger/protocol/external_apps/bid/bid_rpc"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
)

type Service struct {
	balances       *balance.Store
	currencies     *balance.CurrencySet
	ons            *ons.DomainStore
	logger         *log.Logger
	bidMaster  	   *bid_data.BidMasterStore
}
func Name() string {
	return "bid_query"
}

func NewService(balances *balance.Store, currencies *balance.CurrencySet,
	domains *ons.DomainStore, logger *log.Logger, bidMaster *bid_data.BidMasterStore) *Service {
	return &Service{
		currencies:     currencies,
		balances:       balances,
		ons:            domains,
		logger:         logger,
		bidMaster:      bidMaster,
	}
}

func (svc *Service) ShowBidConv(req bid_rpc.ListBidConvRequest, reply *bid_rpc.ListBidConvsReply) error {
	bidConv, _, err := svc.bidMaster.BidConv.QueryAllStores(req.BidConvId)
	if err != nil {
		svc.logger.Error("error getting bid conversation", err)
		return bid_data.ErrGettingBidConv
	}

	bidOffers := svc.bidMaster.BidOffer.GetOffers(bidConv.BidConvId, bid_data.BidOfferInvalid, bid_data.TypeInvalid)

	bcs := bid_rpc.BidConvStat{
		BidConv: *bidConv,
		Offers:  bidOffers,
	}

	*reply = bid_rpc.ListBidConvsReply{
		BidConvStats: []bid_rpc.BidConvStat{bcs},
		Height:       svc.bidMaster.BidConv.GetState().Version(),
	}
	return nil
}

// list single proposal by id or list proposals
func (svc *Service) ListBidConvs(req bid_rpc.ListBidConvsRequest, reply *bid_rpc.ListBidConvsReply) error {
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
	var bidConvs []bid_data.BidConv
	if req.State != bid_data.BidStateInvalid {
		bidConvs = svc.bidMaster.BidConv.FilterBidConvs(req.State, req.Owner, req.AssetName, req.AssetType, req.Bidder)
	} else { // Query in all stores otherwise
		active := svc.bidMaster.BidConv.FilterBidConvs(bid_data.BidStateActive, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		succeed := svc.bidMaster.BidConv.FilterBidConvs(bid_data.BidStateSucceed, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		rejected := svc.bidMaster.BidConv.FilterBidConvs(bid_data.BidStateRejected, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		expired := svc.bidMaster.BidConv.FilterBidConvs(bid_data.BidStateExpired, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		cancelled := svc.bidMaster.BidConv.FilterBidConvs(bid_data.BidStateCancelled, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		bidConvs = append(bidConvs, active...)
		bidConvs = append(bidConvs, succeed...)
		bidConvs = append(bidConvs, rejected...)
		bidConvs = append(bidConvs, expired...)
		bidConvs = append(bidConvs, cancelled...)
	}

	// Organize reply packet:
	// Bid conversations and their offers
	bidConvStats := make([]bid_rpc.BidConvStat, len(bidConvs))
	for i, bidConv := range bidConvs {
		bidOffers := svc.bidMaster.BidOffer.GetOffers(bidConv.BidConvId, bid_data.BidOfferInvalid, bid_data.TypeInvalid)
		bcs := bid_rpc.BidConvStat{
			BidConv: bidConv,
			Offers: bidOffers,
		}
		bidConvStats[i] = bcs
	}

	*reply = bid_rpc.ListBidConvsReply{
		BidConvStats: bidConvStats,
		Height:       svc.bidMaster.BidConv.GetState().Version(),
	}
	return nil
}

// list active offers
func (svc *Service) ListActiveOffers(req bid_rpc.ListActiveOffersRequest, reply *bid_rpc.ListActiveOffersReply) error {
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
	offers := svc.bidMaster.BidOffer.GetOffers("", bid_data.BidOfferActive, req.OfferType)
	activeOfferStats := make([]bid_rpc.ActiveOfferStat, len(offers))
	for i, offer := range offers {
		// get corresponding bid conversation to show the detail
		bidConv, err := svc.bidMaster.BidConv.WithPrefixType(bid_data.BidStateActive).Get(offer.BidConvId)
		if err != nil {
			svc.logger.Error("error getting bid conversation", err)
			return bid_data.ErrGettingBidConv
		}
		if len(req.Bidder) != 0 && !req.Bidder.Equal(bidConv.Bidder) {
			continue
		}
		if len(req.Owner) != 0 && !req.Owner.Equal(bidConv.AssetOwner) {
			continue
		}
		if req.AssetType != bid_data.BidAssetInvalid && req.AssetType != bidConv.AssetType {
			continue
		}
		if req.AssetName != bidConv.AssetName {
			continue
		}
		aos := bid_rpc.ActiveOfferStat{
			ActiveOffer: offer,
			BidConv: *bidConv,
		}
		activeOfferStats[i] = aos
	}

	*reply = bid_rpc.ListActiveOffersReply{
		ActiveOffers: activeOfferStats,
		Height:       svc.bidMaster.BidConv.GetState().Version(),
	}
	return nil
}