/*

 */

package bitcoin

import (
	"errors"

	"github.com/Oneledger/protocol/action"
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

	DefaultLastUpdateHeight = 0
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
	Multisig *keys.MultiSig `json:"multisig"`

	// State tracks the current state of the tracker, Also used for locking distributed access
	State TrackerState `json:"state"`

	// LastUpdateHeight logs the last update height of the tracker
	LastUpdateHeight int64 `json:"lastUpdateHeight"`

	CurrentUTXO *UTXO
	ProcessUTXO *UTXO

	NextLockScript        []byte
	NextLockScriptAddress []byte

	ProcessOwner action.Address
}

func NewTracker(lockScript, lockScriptAddress []byte) *Tracker {

	return &Tracker{
		State:                 AvailableTrackerState,
		LastUpdateHeight:      DefaultLastUpdateHeight,
		CurrentUTXO:           nil,
		NextLockScript:        lockScript,
		NextLockScriptAddress: lockScriptAddress,
	}
}

// GetBalance gets the current balance of the utxo tracker
func (t *Tracker) GetBalance() int64 {
	if t.CurrentUTXO == nil {
		return 0
	}

	return t.CurrentUTXO.Balance
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

	return t.NextLockScriptAddress, nil
}

func (t *Tracker) ProcessLock(newUTXO *UTXO,
	txn []byte, validatorsPubKeys []*btcutil.AddressPubKey,
) error {

	if t.IsBusy() {
		return ErrTrackerBusy
	}

	t.ProcessUTXO = newUTXO
	t.State = BusySigningTrackerState

	t.Multisig = &keys.MultiSig{}

	signers := make([]keys.Address, len(validatorsPubKeys))
	for i := range validatorsPubKeys {
		signers[i] = validatorsPubKeys[i].ScriptAddress()
	}

	threshold := (len(signers) * 2 / 3) + 1
	err := t.Multisig.Init(txn, threshold, signers)

	return err
}

func (t *Tracker) AddSignature(pubKey keys.PublicKey,
	signatureBytes []byte, validatorPubKey *btcutil.AddressPubKey) error {

	if t.State != BusySigningTrackerState {
		return ErrTrackerNotCollectionSignatures
	}

	index, err := t.Multisig.GetSignerIndex(validatorPubKey.ScriptAddress())
	if err != nil {
		return err
	}

	s := keys.Signature{
		Index:  index,
		PubKey: pubKey,
		Signed: signatureBytes,
	}

	return t.Multisig.AddSignature(s)
}

func (t *Tracker) HasEnoughSignatures() (bool, error) {

	if t.State != BusySigningTrackerState {
		return false, ErrTrackerNotCollectionSignatures
	}

	if t.Multisig.IsValid() {
		return true, nil
	}

	return false, nil
}

func (t *Tracker) StateChangeBroadcast() bool {
	if ok, _ := t.HasEnoughSignatures(); ok {
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
		signatures[i] = signed.Signed
	}

	return signatures
}
