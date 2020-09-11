package bid_error

const (
	//Bidding Error
	BidErrInvalidBidConvId                 = 990001
	BidErrInvalidAsset                     = 990002
	BidErrFailedCreateBidConv              = 990003
	BidErrGettingBidMasterStore            = 990004
	BidErrBidConvNotFound                  = 990005
	BidErrGettingBidConv                   = 990006
	BidErrExpiredBid                       = 990007
	BidErrGettingActiveOffer               = 990008
	BidErrGettingActiveBidOffer            = 990009
	BidErrGettingActiveCounterOffer        = 990010
	BidErrDeactivateOffer                  = 990011
	BidErrCloseBidConv                     = 990012
	BidErrAmountMoreThanActiveCounterOffer = 990013
	BidErrAmountLessThanActiveOffer        = 990014
	BidErrLockAmount                       = 990015
	BidErrUnlockAmount                     = 990016
	BidErrAddingOffer                      = 990017
	BidErrAddingCounterOffer               = 990018
	BidErrInvalidDeadline                  = 990019
	BidErrActiveBidConvExists              = 990020
	BidErrAddingBidConvToActiveStore       = 990021
	BidErrWrongBidder                      = 990022
	BidErrWrongAssetOwner                  = 990023
	BidErrDeductingAmountFromBidder        = 990024
	BidErrAdddingAmountToOwner             = 990025
	BidErrSetOffer                         = 990026
	BidErrAssertingBidMasterStore          = 990027
	BidErrAddingBidConvToTargetStore       = 990028
	BidErrDeletingBidConvFromActiveStore   = 990029
	BidErrInvalidOfferType                 = 990030
	BidErrFailedToDeleteActiveOffer        = 990031
	BidErrInvalidBidderDecision            = 990032
	BidErrInvalidOwnerDecision             = 990033
	BidErrFailedToExchangeAsset            = 990034
)
