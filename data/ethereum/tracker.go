package ethereum

import (
	"errors"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type TrackerState int

const (
	AvailableTrackerState = iota
	NewTrackerState
	BusyLockingTrackerState
	BusySigningTrackerState
	BusyBroadcastingTrackerState
	BusyFinalizingTrackerState
	BusyMintingCoin
)

var NilTxHash *chainhash.Hash

var (
	ErrTrackerBusy                    = errors.New("tracker is busy")
	ErrTrackerNotCollectionSignatures = errors.New("tracker not collecting signatures")
)

func init() {
	NilTxHash, _ = chainhash.NewHash([]byte{})
}

// Tracker
type Tracker struct {
	// State tracks the current state of the tracker, Also used for locking distributed access
	State       TrackerState `json:"state"`
	SignedETHTx []byte
	//UserOlAddress common.Address
	ProcessOwner  keys.Address
	FinalityVotes []keys.Address
	TrackerName   ethereum.TrackerName
}

func NewTracker(name common.Hash) *Tracker {

	return &Tracker{
		State:       NewTrackerState,
		TrackerName: name,
	}
}
