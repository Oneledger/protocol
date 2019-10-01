/*

 */

package bitcoin

import (
	"errors"

	"github.com/Oneledger/protocol/serialize"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

var (
	ErrWrongTrackerAdapter = errors.New("wrong tracker adapter")
)

type TrackerData struct {
	// State tracks the current state of the tracker, Also used for locking distributed access
	State TrackerState `json:"state"`

	// LastUpdateHeight logs the last update height of the tracker
	LastUpdateHeight int64 `json:"lastUpdateHeight"`

	// PreviousTxID is the log of the last successful transaction in the tracker
	PreviousTxID *chainhash.Hash

	LatestUTXO *UTXO

	MultiSigData []byte
}

func (TrackerData) SerialTag() string {
	return "data.bitcoin.Tracker"
}

func (t *Tracker) NewDataInstance() serialize.Data {
	return &TrackerData{}
}

func (t *Tracker) Data() serialize.Data {
	return &TrackerData{
		t.State,
		t.LastUpdateHeight,
		t.PreviousTxID,
		t.LatestUTXO,
		t.Multisig.Bytes(),
	}
}

func (t *Tracker) SetData(obj interface{}) error {
	td, ok := obj.(*TrackerData)
	if !ok {
		return ErrWrongTrackerAdapter
	}

	t.State = td.State
	t.LastUpdateHeight = td.LastUpdateHeight
	t.PreviousTxID = td.PreviousTxID
	t.LatestUTXO = td.LatestUTXO
	err := t.Multisig.FromBytes(td.MultiSigData)
	if err != nil {
		return err
	}

	return nil
}
