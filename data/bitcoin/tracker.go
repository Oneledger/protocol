/*

 */

package bitcoin

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type TrackerState int

const (
	AvailableTrackerState = iota
	BusyLockingTrackerState
	BusySigningTrackerState
	BusyBroadcastingTrackerState
	BusyFinalizingTrackerState
)

var NilTxHash *chainhash.Hash

func init() {
	NilTxHash, _ = chainhash.NewHash([]byte{})
}

// Tracker
type Tracker struct {
	// Multisig manages the signature collection and storage in a distributed way
	Multisig keys.MultiSigner `json:"multisig"`

	// State tracks the current state of the tracker, Also used for locking distributed access
	State TrackerState `json:"state"`

	// LastUpdateHeight logs the last update height of the tracker
	LastUpdateHeight int64 `json:"lastUpdateHeight"`

	// PreviousTxID is the log of the last successful transaction in the tracker
	PreviousTxID *chainhash.Hash

	LatestUTXO *UTXO

	NextLockScriptAddress []byte
}

func NewTracker() *Tracker {

	return &Tracker{
		State:            AvailableTrackerState,
		LastUpdateHeight: 0,
		PreviousTxID:     NilTxHash,
		LatestUTXO:       nil,
	}
}

// GetBalance gets the current balance of the utxo tracker
func (t *Tracker) GetBalance() int64 {
	if t.LatestUTXO == nil {
		return 0
	}

	return t.LatestUTXO.Balance
}

// IsAvailable returns whether the tracker is available for new transaction
func (t *Tracker) IsAvailable() bool {
	return t.State == AvailableTrackerState
}

// IsBusy returns whether the tracker is in any of the busy states
func (t *Tracker) IsBusy() bool {
	return t.State != AvailableTrackerState
}
