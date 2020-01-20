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

	ProcessTypeNone    ProcessType = 0x00
	ProcessTypeLock    ProcessType = 0x01
	ProcessTypeRedeem  ProcessType = 0x02
	ProcessTypeLockERC ProcessType = 0x03
	ProcessTypeRedeemERC ProcessType = 0x04
)

var (
	ErrTrackerNotFound    = errors.New("tracker not found")
	errTrackerInvalidVote = errors.New("vote information is invalid")
)

type ProcessType int8

func GetProcessTypeString (t ProcessType) string {
	switch t{
	case 0x00 :return "ProcessTypeNone"
	case 0x01 :return "ETH LOCK"
	case 0x02 :return "ETH REDEEM"
	case 0x03 :return "ERC LOCK"
	case 0x04 :return "ERC REDEEM"
	}
	return "UNKNOWN TYPE"
}

type Vote uint8
