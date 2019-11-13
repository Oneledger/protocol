package ethereum

import (
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

var (
	ErrTrackerNotFound = errors.New("tracker not found")
)

type TrackerStore struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte
}

func (ts *TrackerStore) Get(key ethereum.TrackerName) (Tracker, error) {
	tracker := Tracker{}
	prefixed := append(ts.prefix, key.Bytes()...)
	data, err := ts.state.Get(prefixed)
	if err != nil {
		return tracker, err
	}

	err = ts.szlr.Deserialize(data, tracker)

	return tracker, err
}

func (ts *TrackerStore) Set(tracker Tracker) error {
	prefixed := append(ts.prefix, tracker.TrackerName.Bytes()...)
	data, err := ts.szlr.Serialize(tracker)
	if err != nil {
		return err
	}
	err = ts.state.Set(prefixed, data)

	return err
}

func (ts *TrackerStore) Exists(key ethereum.TrackerName) bool {
	prefixed := append(ts.prefix, key.Bytes()...)
	return ts.state.Exists(prefixed)
}

func (ts *TrackerStore) Delete(key ethereum.TrackerName) (bool, error) {
	prefixed := append(ts.prefix, key.Bytes()...)
	return ts.state.Delete(prefixed)
}

/*
func (ts *TrackerStore) GetIterator() storage.Iteratable {
	panic("implement me")
}

func NewTrackerStore(prefix string, state *storage.State) *TrackerStore {
	return &TrackerStore{
		State:  state,
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
		prefix: storage.Prefix(prefix),
	}
}

// WithState updates the storage state of the tracker and returns the tracker address back
func (ts *TrackerStore) WithState(state *storage.State) *TrackerStore {
	ts.State = state
	return ts
}*/
