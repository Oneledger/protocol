package ethereum

import (
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type TrackerStore struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte
	cdOpt  *ethereum.ChainDriverOption
}

func (ts *TrackerStore) Get(key ethereum.TrackerName) (*Tracker, error) {
	tracker := &Tracker{}
	prefixed := append(ts.prefix, key.Bytes()...)
	data, err := ts.state.Get(prefixed)
	if err != nil {
		return nil, err
	}

	err = ts.szlr.Deserialize(data, tracker)

	return tracker, err
}

func (ts *TrackerStore) Set(tracker *Tracker) error {
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

func (ts *TrackerStore) GetIterable() storage.Iterable {
	return ts.state.GetIterable()
}

func (ts *TrackerStore) Iterate(fn func(name *ethereum.TrackerName, tracker *Tracker) bool) (stopped bool) {
	return ts.state.IterateRange(
		ts.prefix,
		storage.Rangefix(string(ts.prefix)),
		true,
		func(key, value []byte) bool {
			name := &ethereum.TrackerName{}
			name.SetBytes([]byte(string(key[len(ts.prefix):])))

			tracker := &Tracker{}
			err := ts.szlr.Deserialize(value, tracker)
			if err != nil {
				return false
			}
			return fn(name, tracker)
		},
	)
}

func NewTrackerStore(prefix string, state *storage.State) *TrackerStore {
	return &TrackerStore{
		state:  state,
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
		prefix: storage.Prefix(prefix),
		cdOpt:  &ethereum.ChainDriverOption{},
	}
}

// WithState updates the storage state of the tracker and returns the tracker address back
func (ts *TrackerStore) WithState(state *storage.State) *TrackerStore {
	ts.state = state
	return ts
}

func (ts *TrackerStore) SetupOption(opt *ethereum.ChainDriverOption) {
	ts.cdOpt = opt
}

func (ts *TrackerStore) GetOption() *ethereum.ChainDriverOption {
	return ts.cdOpt
}
