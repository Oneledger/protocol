package bid_rpc_query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
	"github.com/Oneledger/protocol/external_apps/bid/bid_rpc"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/pkg/errors"
)

type Service struct {
	ext            client.ExtServiceContext
	balances       *balance.Store
	currencies     *balance.CurrencySet
	validators     *identity.ValidatorStore
	ons            *ons.DomainStore
	feePool        *fees.Store
	governance     *governance.Store
	logger         *log.Logger
	extStores  	   data.Router
}
//todo delete this if not used
func Name() string {
	return "ext_query"
}

func NewService(ctx client.ExtServiceContext, balances *balance.Store, currencies *balance.CurrencySet, validators *identity.ValidatorStore,
	domains *ons.DomainStore, govern *governance.Store, feePool *fees.Store, logger *log.Logger, extStores data.Router) *Service {
	return &Service{
		ext:            ctx,
		currencies:     currencies,
		balances:       balances,
		validators:     validators,
		ons:            domains,
		feePool:        feePool,
		logger:         logger,
		governance:     govern,
		extStores:      extStores,
	}
}

func (svc *Service) ShowBidConv(req bid_rpc.ListBidConvRequest, reply *bid_rpc.ListBidConvsReply) error {
	bidMaster, err := GetBidMasterStore(svc.extStores)
	if err != nil {
		return err
	}
	bidConv, _, err := bidMaster.BidConv.QueryAllStores(req.BidConvId)
	if err != nil {
		svc.logger.Error("error getting bid conversation", err)
		return codes.ErrGettingBidConv
	}

	bidOffers := bidMaster.BidOffer.GetOffers(bidConv.BidConvId, bid_data.BidOfferInvalid, bid_data.TypeInvalid)

	bcs := bid_rpc.BidConvStat{
		BidConv: *bidConv,
		Offers:  bidOffers,
	}

	*reply = bid_rpc.ListBidConvsReply{
		BidConvStats: []bid_rpc.BidConvStat{bcs},
		Height:       bidMaster.BidConv.GetState().Version(),
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

	bidMaster, err := GetBidMasterStore(svc.extStores)
	if err != nil {
		return err
	}

	// Query in single store if specified
	var bidConvs []bid_data.BidConv
	if req.State != bid_data.BidStateInvalid {
		bidConvs = bidMaster.BidConv.FilterBidConvs(req.State, req.Owner, req.AssetName, req.AssetType, req.Bidder)
	} else { // Query in all stores otherwise
		active := bidMaster.BidConv.FilterBidConvs(bid_data.BidStateActive, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		succeed := bidMaster.BidConv.FilterBidConvs(bid_data.BidStateSucceed, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		rejected := bidMaster.BidConv.FilterBidConvs(bid_data.BidStateRejected, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		expired := bidMaster.BidConv.FilterBidConvs(bid_data.BidStateExpired, req.Owner, req.AssetName, req.AssetType, req.Bidder)
		cancelled := bidMaster.BidConv.FilterBidConvs(bid_data.BidStateCancelled, req.Owner, req.AssetName, req.AssetType, req.Bidder)
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
		bidOffers := bidMaster.BidOffer.GetOffers(bidConv.BidConvId, bid_data.BidOfferInvalid, bid_data.TypeInvalid)
		bcs := bid_rpc.BidConvStat{
			BidConv: bidConv,
			Offers: bidOffers,
		}
		bidConvStats[i] = bcs
	}

	*reply = bid_rpc.ListBidConvsReply{
		BidConvStats: bidConvStats,
		Height:       bidMaster.BidConv.GetState().Version(),
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
	bidMaster, err := GetBidMasterStore(svc.extStores)
	if err != nil {
		return err
	}

	// get all active offers
	offers := bidMaster.BidOffer.GetOffers("", bid_data.BidOfferActive, req.OfferType)
	activeOfferStats := make([]bid_rpc.ActiveOfferStat, len(offers))
	for i, offer := range offers {
		// get corresponding bid conversation to show the detail
		bidConv, err := bidMaster.BidConv.WithPrefixType(bid_data.BidStateActive).Get(offer.BidConvId)
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
		if req.AssetType != bid_data.BidAssetInvalid && req.AssetType != bidConv.AssetType {
			continue
		}
		if req.AssetName != bidConv.Asset.ToString() {
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
		Height:       bidMaster.BidConv.GetState().Version(),
	}
	return nil
}

func GetBidMasterStore(extStores data.Router) (*bid_data.BidMasterStore, error) {
	store, err := extStores.Get("bidMaster")
	if err != nil {
		return nil, err
	}
	bidMasterStore, ok := store.(*bid_data.BidMasterStore)
	if ok == false {
		return nil, err
	}

	return bidMasterStore, nil
}