package ethereum

const (
	New TrackerState = iota
	BusyBroadcasting
	BusyFinalizing
	Finalized
	Released

	BROADCASTING    string = "broadcasting"
	FINALIZING      string = "finalizing"
	FINALIZE        string = "finalize"
	MINTING         string = "minting"
	CLEANUP         string = "cleanup"
	SIGNING         string = "signing"
	VERIFYREDEEM    string = "verifyredeem"
	BURN            string = "burn"

	ProcessTypeNone   ProcessType = 0x00
	ProcessTypeLock   ProcessType = 0x01
	ProcessTypeRedeem ProcessType = 0x02
)

type ProcessType int8

type Vote uint8
