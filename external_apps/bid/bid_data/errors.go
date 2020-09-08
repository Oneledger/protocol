package bid_data

import (
	"github.com/Oneledger/protocol/external_apps/bid/bid_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrInvalidBidConvId                 = codes.ProtocolError{bid_error.BidErrInvalidBidConvId, "invalid bid conversation id"}
	ErrInvalidAsset                     = codes.ProtocolError{bid_error.BidErrInvalidAsset, "invalid bid asset"}
	ErrFailedCreateBidConv              = codes.ProtocolError{bid_error.BidErrFailedCreateBidConv, "failed to create bid conversation"}
	ErrGettingBidMasterStore            = codes.ProtocolError{bid_error.BidErrGettingBidMasterStore, "failed to get bid master store"}
	ErrBidConvNotFound                  = codes.ProtocolError{bid_error.BidErrBidConvNotFound, "bid conversation not found"}
	ErrGettingBidConv                   = codes.ProtocolError{bid_error.BidErrGettingBidConv, "failed to get bid conversation"}
	ErrExpiredBid                       = codes.ProtocolError{bid_error.BidErrExpiredBid, "the bid is already expired"}
	ErrGettingActiveOffer               = codes.ProtocolError{bid_error.BidErrGettingActiveOffer, "failed to get active offer"}
	ErrTooManyActiveOffers              = codes.ProtocolError{bid_error.BidErrTooManyActiveOffers, "error there should be only one active offer"}
	ErrGettingActiveBidOffer            = codes.ProtocolError{bid_error.BidErrGettingActiveBidOffer, "failed to get active bid offer"}
	ErrGettingActiveCounterOffer        = codes.ProtocolError{bid_error.BidErrGettingActiveCounterOffer, "failed to get active counter offer"}
	ErrDeactivateOffer                  = codes.ProtocolError{bid_error.BidErrDeactivateOffer, "failed to deactivate offer"}
	ErrCloseBidConv                     = codes.ProtocolError{bid_error.BidErrCloseBidConv, "failed to close bid conversation"}
	ErrActiveCounterOfferExists         = codes.ProtocolError{bid_error.BidErrActiveCounterOfferExists, "there is active counter offer"}
	ErrActiveBidOfferExists             = codes.ProtocolError{bid_error.BidErrActiveBidOfferExists, "there is active bid offer"}
	ErrAmountMoreThanActiveCounterOffer = codes.ProtocolError{bid_error.BidErrAmountMoreThanActiveCounterOffer, "amount should not be larger than active counter offer amount"}
	ErrAmountLessThanActiveOffer        = codes.ProtocolError{bid_error.BidErrAmountLessThanActiveOffer, "amount should not be less than active bid offer amount"}
	ErrLockAmount                       = codes.ProtocolError{bid_error.BidErrLockAmount, "failed to lock amount"}
	ErrUnlockAmount                     = codes.ProtocolError{bid_error.BidErrUnlockAmount, "failed to unlock amount"}
	ErrAddingOffer                      = codes.ProtocolError{bid_error.BidErrAddingOffer, "failed to add bid offer"}
	ErrAddingCounterOffer               = codes.ProtocolError{bid_error.BidErrAddingCounterOffer, "failed to add counter offer"}
	ErrInvalidDeadline                  = codes.ProtocolError{bid_error.BidErrInvalidDeadline, "invalid bid conversation deadline"}
	ErrActiveBidConvExists              = codes.ProtocolError{bid_error.BidErrActiveBidConvExists, "bid conversation with same id already exists"}
	ErrAddingBidConvToActiveStore       = codes.ProtocolError{bid_error.BidErrAddingBidConvToActiveStore, "failed to add bid conversation to active store"}
	ErrWrongBidder                      = codes.ProtocolError{bid_error.BidErrWrongBidder, "bidder not match"}
	ErrWrongAssetOwner                  = codes.ProtocolError{bid_error.BidErrWrongAssetOwner, "asset owner not match"}
	ErrDeductingAmountFromBidder        = codes.ProtocolError{bid_error.BidErrDeductingAmountFromBidder, "failed to deduct amount from bidder"}
	ErrAdddingAmountToOwner             = codes.ProtocolError{bid_error.BidErrAdddingAmountToOwner, "failed to add amount to asset owner"}
	ErrUpdateOffer                      = codes.ProtocolError{Code: bid_error.BidErrUpdateOffer, Msg: "failed to update offer"}
	ErrAssertingBidMasterStore          = codes.ProtocolError{Code: bid_error.BidErrAssertingBidMasterStore, Msg: "failed to assert bid master store"}
	ErrAddingBidConvToTargetStore       = codes.ProtocolError{Code: bid_error.BidErrAddingBidConvToTargetStore, Msg: "failed to add bid conversation to target store"}
	ErrDeletingBidConvFromActiveStore   = codes.ProtocolError{Code: bid_error.BidErrDeletingBidConvFromActiveStore, Msg: "failed to delete bid conversation from active store"}
	ErrInvalidOfferType                 = codes.ProtocolError{Code: bid_error.BidErrInvalidOfferType, Msg: "invalid offer type"}
)
