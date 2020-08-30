package ethereum

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	// Options Objects from store
	ErrETHTrackerExists      = codes.ProtocolError{codes.ETHTrackerExists, "Tracker Already exists"}
	ErrETHTrackerUnableToSet = codes.ProtocolError{codes.ETHTrackerUnabletoSet, "Unable to set ETH tracker"}
)
