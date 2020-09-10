package bid_error

const (
	//Bidding Error
	BidErrInvalidBidConvId = 990001
	BidErrInvalidAsset = 990002
	BidErrFailedCreateBidConv = 990003
	BidErrGettingBidMasterStore = 990004
	BidErrBidConvNotFound = 990005
	BidErrGettingBidConv = 990006
	BidErrExpiredBid = 990007
	BidErrGettingActiveOffer = 990008
	BidErrTooManyActiveOffers = 990009
	BidErrGettingActiveBidOffer = 990010
	BidErrGettingActiveCounterOffer = 990011
	BidErrDeactivateOffer = 990012
	BidErrCloseBidConv = 990013
	BidErrActiveCounterOfferExists = 990014
	BidErrActiveBidOfferExists = 990015
	BidErrAmountMoreThanActiveCounterOffer = 990016
	BidErrAmountLessThanActiveOffer = 990017
	BidErrLockAmount = 990018
	BidErrUnlockAmount                   = 990019
	BidErrAddingOffer                    = 990020
	BidErrAddingCounterOffer             = 990021
	BidErrInvalidDeadline                = 990022
	BidErrActiveBidConvExists            = 990023
	BidErrAddingBidConvToActiveStore     = 990024
	BidErrWrongBidder                    = 990025
	BidErrWrongAssetOwner                = 990026
	BidErrDeductingAmountFromBidder      = 990027
	BidErrAdddingAmountToOwner           = 990028
	BidErrSetOffer                       = 990029
	BidErrAssertingBidMasterStore        = 990030
	BidErrAddingBidConvToTargetStore     = 990031
	BidErrDeletingBidConvFromActiveStore = 990032
	BidErrInvalidOfferType               = 990033
	BidErrFailedToDeleteActiveOffer      = 990034
	BidErrInvalidBidderDecision			 = 990035
	BidErrInvalidOwnerDecision			 = 990036
	BidErrFailedToExchangeAsset			 = 990037
)
