package storage

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

func (s *State) GetIterator() *Iterator {
	return s.current.GetIterator()
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
