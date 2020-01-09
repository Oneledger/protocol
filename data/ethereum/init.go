package ethereum

import "errors"

const (
	New TrackerState = iota
	BusyBroadcasting
	BroadcastSuccess
	BusyFinalizing
	Finalized
	Released

	BROADCASTING  string = "broadcasting"
	FINALIZING    string = "finalizing"
	FINALIZE      string = "finalize"
	MINTING       string = "minting"
	CLEANUP       string = "cleanup"
	SIGNING       string = "signing"
	VERIFYREDEEM  string = "verifyredeem"
	REDEEMCONFIRM string = "redeemconfirm"
	BURN          string = "burn"

	ProcessTypeNone   ProcessType = 0x00
	ProcessTypeLock   ProcessType = 0x01
	ProcessTypeRedeem ProcessType = 0x02
	ProcessTypeLockERC ProcessType = 0x03
)

var (
	ErrTrackerNotFound    = errors.New("tracker not found")
	errTrackerInvalidVote = errors.New("vote information is invalid")
)

type ProcessType int8

type Vote uint8
