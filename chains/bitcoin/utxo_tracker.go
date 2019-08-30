/*

 */

package bitcoin

import "github.com/btcsuite/btcd/chaincfg/chainhash"

type TrackerState int

const (
	StatusAvailable TrackerState = iota
	StatusBusy
)

type utxoTracker struct {
	UTXO
	PreviousTxID chainhash.Hash
	State        TrackerState
}

// GetBalance gets the current balance of the utxo tracker
func (u *utxoTracker) GetBalance() int64 {
	return u.Balance
}

// IsAvailable returns true if the tracker is available for transactions
func (u *utxoTracker) IsAvailable() bool {
	return u.State == StatusAvailable
}

// IsBusy returns true if the tracker
func (u *utxoTracker) IsBusy() bool {
	return u.State == StatusBusy
}
