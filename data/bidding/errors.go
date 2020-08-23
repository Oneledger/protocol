package bidding

import codes "github.com/Oneledger/protocol/status_codes"

var (
	ErrInvalidBidConvId                 = codes.ProtocolError{codes.BidErrInvalidBidConvId, "invalid bid conversation id"}
	ErrInvalidAsset                     = codes.ProtocolError{codes.BidErrInvalidAsset, "invalid bid asset"}
	ErrFailedCreateBidConv              = codes.ProtocolError{codes.BidErrFailedCreateBidConv, "failed to create bid conversation"}
	ErrGettingBidMasterStore            = codes.ProtocolError{codes.BidErrGettingBidMasterStore, "failed to get bid master store"}
	ErrBidConvNotFound                  = codes.ProtocolError{codes.BidErrBidConvNotFound, "bid conversation not found"}
	ErrGettingBidConv                   = codes.ProtocolError{codes.BidErrGettingBidConv, "failed to get bid conversation"}
	ErrExpiredBid                       = codes.ProtocolError{codes.BidErrExpiredBid, "the bid is already expired"}
	ErrGettingActiveOffer               = codes.ProtocolError{codes.BidErrGettingActiveOffer, "failed to get active offer"}
	ErrTooManyActiveOffers              = codes.ProtocolError{codes.BidErrTooManyActiveOffers, "error there should be only one active offer"}
	ErrGettingActiveBidOffer            = codes.ProtocolError{codes.BidErrGettingActiveBidOffer, "failed to get active bid offer"}
	ErrGettingActiveCounterOffer        = codes.ProtocolError{codes.BidErrGettingActiveCounterOffer, "failed to get active counter offer"}
	ErrDeactivateOffer                  = codes.ProtocolError{codes.BidErrDeactivateOffer, "failed to deactivate offer"}
	ErrCloseBidConv                     = codes.ProtocolError{codes.BidErrCloseBidConv, "failed to close bid conversation"}
	ErrActiveCounterOfferExists         = codes.ProtocolError{codes.BidErrActiveCounterOfferExists, "there is active counter offer"}
	ErrActiveBidOfferExists             = codes.ProtocolError{codes.BidErrActiveBidOfferExists, "there is active bid offer"}
	ErrAmountMoreThanActiveCounterOffer = codes.ProtocolError{codes.BidErrAmountMoreThanActiveCounterOffer, "amount should not be larger than active counter offer amount"}
	ErrAmountLessThanActiveOffer        = codes.ProtocolError{codes.BidErrAmountLessThanActiveOffer, "amount should not be less than active bid offer amount"}
	ErrLockAmount                       = codes.ProtocolError{codes.BidErrLockAmount, "failed to lock amount"}
	ErrUnlockAmount                     = codes.ProtocolError{codes.BidErrUnlockAmount, "failed to unlock amount"}
	ErrAddingOffer                      = codes.ProtocolError{codes.BidErrAddingOffer, "failed to add bid offer"}
	ErrAddingCounterOffer               = codes.ProtocolError{codes.BidErrAddingCounterOffer, "failed to add counter offer"}
	ErrInvalidDeadline                  = codes.ProtocolError{codes.BidErrInvalidDeadline, "invalid bid conversation deadline"}
	ErrActiveBidConvExists              = codes.ProtocolError{codes.BidErrActiveBidConvExists, "bid conversation with same id already exists"}
	ErrAddingBidConvToActiveStore       = codes.ProtocolError{codes.BidErrAddingBidConvToActiveStore, "failed to add bid conversation to active store"}
	ErrWrongBidder                      = codes.ProtocolError{codes.BidErrWrongBidder, "bidder not match"}
	ErrWrongAssetOwner                  = codes.ProtocolError{codes.BidErrWrongAssetOwner, "asset owner not match"}
	ErrDeductingAmountFromBidder        = codes.ProtocolError{codes.BidErrDeductingAmountFromBidder, "failed to deduct amount from bidder"}
	ErrAdddingAmountToOwner             = codes.ProtocolError{codes.BidErrAdddingAmountToOwner, "failed to add amount to asset owner"}
	ErrUpdateOffer                      = codes.ProtocolError{Code: codes.BidErrUpdateOffer, Msg: "failed to update offer"}
	ErrAssertingBidMasterStore          = codes.ProtocolError{Code: codes.BidErrAssertingBidMasterStore, Msg: "failed to assert bid master store"}
	ErrAddingBidConvToTargetStore       = codes.ProtocolError{Code: codes.BidErrAddingBidConvToTargetStore, Msg: "failed to add bid conversation to target store"}
	ErrDeletingBidConvFromActiveStore   = codes.ProtocolError{Code: codes.BidErrDeletingBidConvFromActiveStore, Msg: "failed to delete bid conversation from active store"}
	ErrInvalidOfferType                 = codes.ProtocolError{Code: codes.BidErrInvalidOfferType, Msg: "invalid offer type"}
)
