/*

 */

package bitcoin

import (
	"bytes"
	"strconv"

	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils/transition"

	"github.com/pkg/errors"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/Oneledger/protocol/data/keys"
)

type TrackerState int

const (
	ProcessTypeNone   = 0x00
	ProcessTypeLock   = 0x01
	ProcessTypeRedeem = 0x02
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
	Name string

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
	ProcessType  int

	FinalityVotes []keys.Address
	ResetVotes    []keys.Address
}

func NewTracker(lockScriptAddress []byte, m int, signers []keys.Address) (*Tracker, error) {

	btcMultisig, err := keys.NewBTCMultiSig(nil, m, signers)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing multisig")
	}

	return &Tracker{
		State:                    Available,
		CurrentTxId:              nil,
		CurrentLockScriptAddress: nil,

		ProcessLockScriptAddress: lockScriptAddress,
		Multisig:                 btcMultisig,
	}, nil
}

// GetBalance gets the current balance of the utxo tracker
func (t *Tracker) GetBalanceSatoshi() int64 {

	return t.CurrentBalance
}

// IsAvailable returns whether the tracker is available for new transaction
func (t *Tracker) IsAvailable() bool {
	return t.State == Available
}

// IsBusy returns whether the tracker is in any of the busy states
func (t *Tracker) IsBusy() bool {
	return t.State != Available
}

func (t *Tracker) GetAddress() ([]byte, error) {
	if t.IsBusy() {
		return nil, ErrTrackerBusy
	}

	return t.ProcessLockScriptAddress, nil
}

func (t *Tracker) ProcessLock(newUTXO *UTXO,
	txn []byte, validatorsPubKeys []keys.Address,
) error {

	if t.IsBusy() {
		return ErrTrackerBusy
	}

	t.ProcessBalance = newUTXO.Balance
	t.ProcessUnsignedTx = txn

	t.State = BusySigning

	threshold := (len(validatorsPubKeys) * 2 / 3) + 1

	ms, err := keys.NewBTCMultiSig(txn, threshold, validatorsPubKeys)
	t.Multisig = ms

	return err
}

func (t *Tracker) AddSignature(signatureBytes []byte, addr keys.Address) error {

	if t.State != BusySigning {
		return ErrTrackerNotCollectionSignatures
	}

	index, err := t.Multisig.GetSignerIndex(addr)
	if err != nil {
		return err
	}

	s := keys.BTCSignature{
		Index:   index,
		Address: addr,
		Sign:    signatureBytes,
	}

	return t.Multisig.AddSignature(&s)
}

func (t *Tracker) HasEnoughSignatures() bool {

	if t.State != BusySigning {
		return false
	}

	if t.Multisig.IsValid() {
		return true
	}

	return false
}

func (t *Tracker) StateChangeBroadcast() bool {
	if ok := t.HasEnoughSignatures(); ok {
		t.State = BusyBroadcasting

		return true
	}

	return false
}

func (t *Tracker) GetSignatures() [][]byte {
	if t.State != BusyBroadcasting {
		return nil
	}

	return t.Multisig.GetSignaturesInOrder()
}

func (t *Tracker) GetBalance() int64 {
	return t.CurrentBalance
}

func (t Tracker) NextStep() string {

	switch t.State {
	case Requested:
		return RESERVE
	case BusyScheduleBroadcasting:
		return FREEZE_FOR_BROADCAST
	case BusyScheduleFinalizing:
		return REPORT_BROADCAST
	case Finalized:
		return CLEANUP
	}
	return transition.NOOP

}

func (t *Tracker) GetJobID(state TrackerState) string {
	return t.Name + storage.DB_PREFIX + strconv.Itoa(int(state))
}

func (t *Tracker) HasVotedFinality(addr keys.Address) bool {
	for i := range t.FinalityVotes {
		if bytes.Equal(addr, t.FinalityVotes[i]) {
			return true
		}
	}

	return false
}
