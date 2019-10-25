/*

 */

package bitcoin

import (
	"errors"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
)

type TrackerState int

const (
	AvailableTrackerState = iota
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
	// Multisig manages the signature collection and storage in a distributed way
	Multisig *keys.BTCMultiSig `json:"multisig"`

	// State tracks the current state of the tracker, Also used for locking distributed access
	State TrackerState `json:"state"`

	CurrentTxId              *chainhash.Hash
	CurrentBalance           int64
	CurrentLockScriptAddress []byte

	ProcessTxId              *chainhash.Hash
	ProcessBalance           int64
	ProcessLockScriptAddress []byte
	ProcessUnsignedTx        []byte

	ProcessOwner keys.Address

	FinalityVotes []keys.Address
}

func NewTracker(lockScript, lockScriptAddress []byte) *Tracker {

	return &Tracker{
		State:                    AvailableTrackerState,
		CurrentTxId:              nil,
		CurrentLockScriptAddress: nil,

		ProcessLockScriptAddress: lockScriptAddress,
	}
}

// GetBalance gets the current balance of the utxo tracker
func (t *Tracker) GetBalanceSatoshi() int64 {

	return t.CurrentBalance
}

// IsAvailable returns whether the tracker is available for new transaction
func (t *Tracker) IsAvailable() bool {
	return t.State == AvailableTrackerState
}

// IsBusy returns whether the tracker is in any of the busy states
func (t *Tracker) IsBusy() bool {
	return t.State != AvailableTrackerState
}

func (t *Tracker) GetAddress() ([]byte, error) {
	if t.IsBusy() {
		return nil, ErrTrackerBusy
	}

	return t.ProcessLockScriptAddress, nil
}

func (t *Tracker) ProcessLock(newUTXO *UTXO,
	txn []byte, validatorsPubKeys []btcutil.AddressPubKey,
) error {

	if t.IsBusy() {
		return ErrTrackerBusy
	}

	t.ProcessBalance = newUTXO.Balance
	t.ProcessUnsignedTx = txn

	t.State = BusySigningTrackerState

	threshold := (len(validatorsPubKeys) * 2 / 3) + 1

	ms, err := keys.NewBTCMultiSig(txn, threshold, validatorsPubKeys)
	t.Multisig = ms

	return err
}

func (t *Tracker) AddSignature(signatureBytes []byte,
	validatorPubKey btcutil.AddressPubKey, params *chaincfg.Params) error {

	if t.State != BusySigningTrackerState {
		return ErrTrackerNotCollectionSignatures
	}

	index, err := t.Multisig.GetSignerIndex(validatorPubKey)
	if err != nil {
		return err
	}

	s := keys.BTCSignature{
		Index:  index,
		PubKey: validatorPubKey.PubKey().SerializeCompressed(),
		Sign:   signatureBytes,
	}

	return t.Multisig.AddSignature(&s, params)
}

func (t *Tracker) HasEnoughSignatures() bool {

	if t.State != BusySigningTrackerState {
		return false
	}

	if t.Multisig.IsValid() {
		return true
	}

	return false
}

func (t *Tracker) StateChangeBroadcast() bool {
	if ok := t.HasEnoughSignatures(); ok {
		t.State = BusyBroadcastingTrackerState

		return true
	}

	return false
}

func (t *Tracker) GetSignatures() [][]byte {
	if t.State != BusyBroadcastingTrackerState {
		return nil
	}

	signatures := make([][]byte, 0, len(t.Multisig.Signatures))
	for i, signed := range t.Multisig.GetSignatures() {
		signatures[i] = signed.Sign
	}

	return signatures
}

func (t *Tracker) GetBalance() int64 {
	return t.CurrentBalance
}
