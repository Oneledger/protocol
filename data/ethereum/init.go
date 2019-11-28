package ethereum

const (
	New TrackerState = iota
	BusyBroadcasting
	BusyFinalizing
	Finalized
	Minted

	BROADCASTING string = "broadcasting"
	FINALIZING   string = "finalizing"
	FINALIZE     string = "finalize"
	MINTING      string = "minting"
	CLEANUP      string = "cleanup"

	ProcessTypeNone   ProcessType = 0x00
	ProcessTypeLock   ProcessType = 0x01
	ProcessTypeRedeem ProcessType = 0x02
)

type ProcessType int8

type Vote uint8
