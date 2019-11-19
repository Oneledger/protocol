package ethereum

const (
	New TrackerState = iota
	BusyBroadcasting
	BusyFinalizing
	Finalized
	Minted

	votesThreshold float32 = 0.6667

	BROADCASTING string = "broadcasting"
	FINALIZING   string = "finalizing"
	FINALIZE     string = "finalize"
	MINTING      string = "minting"
	CLEANUP      string = "cleanup"
)
