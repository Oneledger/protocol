package storage

var _ Store = &State{}
var _ Iteratable = &State{}

type State struct {
	cs      *ChainState
	current Store
}

func (s *State) Get(key StoreKey) ([]byte, error) {
	return s.current.Get(key)
}

func (s *State) Set(key StoreKey, value []byte) error {
	return s.current.Set(key, value)
}

func (s *State) Exists(key StoreKey) bool {
	return s.current.Exists(key)
}

func (s *State) Delete(key StoreKey) (bool, error) {
	return s.current.Delete(key)
}

func (s *State) GetIterator() Iteratable {
	return s.current.GetIterator()
}

func (s *State) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
	return s.current.GetIterator().Iterate(fn)
}

func (s *State) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool) {
	return s.current.GetIterator().IterateRange(start, end, ascending, fn)
}

func NewState(state *ChainState) *State {
	return &State{
		cs:      state,
		current: state,
	}
}

func (s *State) WithGas(gc GasCalculator) *State {
	gcs := NewGasChainState(s.cs, gc)
	return &State{
		cs:      s.cs,
		current: gcs,
	}
}

func (s *State) WithoutGas() *State {
	s.current = s.cs
	return s
}

func (s State) Version() int64 {
	return s.cs.Version
}

func (s State) RootHash() []byte {
	return s.cs.Hash
}
