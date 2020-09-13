package bid_rpc

import (
	"github.com/Oneledger/protocol/external_apps/bid/bid_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrGettingBidConvInQuery       = codes.ProtocolError{bid_error.BidErrGettingBidConvInQuery, "failed to get bid conversation in query"}
	ErrGettingActiveOfferInQuery   = codes.ProtocolError{bid_error.BidErrGettingActiveOfferInQuery, "failed to get active offer in query"}
	ErrInvalidOwnerAddressInQuery  = codes.ProtocolError{bid_error.BidErrInvalidOwnerAddressInQuery, "invalid owner address"}
	ErrInvalidBidderAddressInQuery = codes.ProtocolError{bid_error.BidErrInvalidBidderAddressInQuery, "invalid bidder address"}
)
