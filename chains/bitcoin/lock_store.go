/*

 */

package bitcoin

import (
	"fmt"

	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type LockStore struct {
	*storage.ChainState
	serialize.Serializer
}

func NewLockStore(name, dbDir, configDB string, typ storage.StorageType) *LockStore {
	cs := storage.NewChainState(name, dbDir, configDB, typ)

	return &LockStore{
		cs,
		serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (ls *LockStore) GetLatestUTXO(lockName string) (UTXO, error) {

	tracker := utxoTracker{}
	d := ls.Get(storage.StoreKey(lockName), false)
	if len(d) == 0 {
		return tracker.UTXO, errors.New("tracker not found")
	}

	err := ls.Deserialize(d, &tracker)

	return tracker.UTXO, err
}

func (ls *LockStore) UpdateUTXO(lockName string, utxo UTXO) error {
	d := ls.Get(storage.StoreKey(lockName), false)
	if len(d) == 0 {
		d = ls.Get(storage.StoreKey(lockName), true)
	}

	fmt.Println(string(d))
	tracker := utxoTracker{}
	err := ls.Deserialize(d, &tracker)
	if err != nil {
		return err
	}

	if tracker.IsBusy() {
		return errors.New("tracker busy")
	}

	tracker.PreviousTxID = tracker.UTXO.TxID
	tracker.State = StatusBusy
	tracker.UTXO = utxo

	d, err = ls.Serialize(tracker)
	if err != nil {
		return err
	}

	err = ls.Set(storage.StoreKey(lockName), d)
	return err
}

func (ls *LockStore) InitializeTracker(lockName string, utxo UTXO) error {

	doesExist := ls.Exists(storage.StoreKey(lockName))
	if doesExist {
		return errors.New("already initialized")
	}
	doesExist = ls.ExistsUncommitted(storage.StoreKey(lockName))
	if doesExist {
		return errors.New("already initialized")
	}

	tracker := utxoTracker{}

	tracker.PreviousTxID = tracker.UTXO.TxID
	tracker.State = StatusBusy
	tracker.UTXO = utxo

	d, err := ls.Serialize(tracker)
	if err != nil {
		return err
	}

	err = ls.Set(storage.StoreKey(lockName), d)
	return err
}

func (ls *LockStore) FinalizeTracker(lockName string) error {
	d := ls.Get(storage.StoreKey(lockName), false)
	if len(d) == 0 {
		return errors.New("tracker not found")
	}

	tracker := utxoTracker{}
	err := ls.Deserialize(d, &tracker)
	if err != nil {
		return err
	}

	if tracker.IsAvailable() {
		return errors.New("tracker not busy")
	}

	tracker.State = StatusAvailable
	d, err = ls.Serialize(tracker)
	if err != nil {
		return err
	}

	err = ls.Set(storage.StoreKey(lockName), d)
	return err
}

func (ls *LockStore) GetAllTrackers() map[string]utxoTracker {
	trackerData := ls.FindAll()

	var trackerMap = map[string]utxoTracker{}
	for name, data := range trackerData {

		tracker := utxoTracker{}
		err := ls.Deserialize(data, &tracker)
		if err == nil {
			trackerMap[name] = tracker
		}
	}

	return trackerMap
}

func (ls *LockStore) GetActiveTrackers() map[string]utxoTracker {
	trackerData := ls.FindAll()

	var trackerMap = map[string]utxoTracker{}
	for name, data := range trackerData {

		tracker := utxoTracker{}
		err := ls.Deserialize(data, &tracker)
		if err == nil && tracker.IsAvailable() {
			trackerMap[name] = tracker
		}
	}

	return trackerMap
}
