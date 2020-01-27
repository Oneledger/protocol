package storage

var _ Store = &State{}
var _ Iteratable = &State{}

type State struct {
	cs     *ChainState
	cache  Store
	gc     GasCalculator
	delete Store
}

func (s *State) Get(key StoreKey) ([]byte, error) {
	// Get the cache first
	result, err := s.cache.Get(key)
	if err == nil {
		// if got result, return directly
		return result, err
	}
	// if didn't get result in cache, get from ChainState
	return s.cs.Get(key)
}

func (s *State) Set(key StoreKey, value []byte) error {
	// set only for cache, waiting to be committed
	return s.cache.Set(key, value)
}

func (s *State) Exists(key StoreKey) bool {
	// check existence in cache, because it's cheaper
	exist := s.cache.Exists(key)
	if !exist {
		// if not existed in cache, check ChainState
		return s.cs.Exists(key)
	}
	return exist
}

func (s *State) Delete(key StoreKey) (bool, error) {
	//cache delete is always true
	_, _ = s.cache.Delete(key)
	err := s.delete.Set(key, []byte{127})
	if err != nil {
		return false, err
	}
	return true, nil
}

// This only Iterate for the ChainState
func (s *State) GetIterator() Iteratable {
	return s
}

func (s *State) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
	keys := make([]StoreKey, 0, 100)
	s.cs.Iterate(func(key, value []byte) bool {
		keys = append(keys, key)
		return false
	})

	for _, key := range keys {
		value, err := s.Get(key)
		if err != nil {
			continue
		}
		stop := fn(key, value)
		if stop {
			return true
		}
	}
	return true
}

func (s *State) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool) {
	keys := make([]StoreKey, 0, 100)
	s.cs.IterateRange(start, end, ascending, func(key, value []byte) bool {
		keys = append(keys, key)
		return false
	})

	for _, key := range keys {
		value, err := s.Get(key)
		if err != nil {
			continue
		}
		stop := fn(key, value)
		if stop {
			return true
		}
	}
	return true
}

func NewState(state *ChainState) *State {
	return &State{
		cs:     state,
		cache:  NewStorage(CACHE, "state"),
		gc:     NewGasCalculator(0),
		delete: NewStorage(CACHE, "state_delete"),
	}
}

func (s *State) WithGas(gc GasCalculator) *State {
	gs := NewGasStore(s.cache, gc)
	del := NewGasStore(s.delete, gc)
	return &State{
		cs:     s.cs,
		cache:  gs,
		gc:     gc,
		delete: del,
	}
}

func (s *State) WithoutGas() *State {
	s.cache = NewStorage(CACHE, "state")
	//s.gc = NewGasCalculator(0)
	s.delete = NewStorage(CACHE, "state_delete")
	return s
}

func (s State) Version() int64 {
	return s.cs.Version
}

func (s State) RootHash() []byte {
	return s.cs.Hash
}

func (s State) Write() bool {
	s.cache.GetIterator().Iterate(func(key []byte, value []byte) bool {
		_ = s.cs.Set(key, value)
		return false
	})
	s.delete.GetIterator().Iterate(func(key, value []byte) bool {
		_, _ = s.cs.Delete(key)
		return false
	})
	return true
}

func (s *State) Commit() (hash []byte, version int64) {
	s.Write()
	s.cache = NewStorage(CACHE, "state")
	return s.cs.Commit()
}

func (s *State) ConsumedGas() Gas {
	return s.gc.GetConsumed()
}

func (s *State) ConsumeUpfront(gas Gas) bool {
	return s.gc.Consume(gas, FLAT, true)
}

func (s *State) ConsumeVerifySigGas(gas Gas) bool {
	return s.gc.Consume(gas, VERIFYSIG, true)
}

func (s *State) ConsumeStorageGas(gas Gas) bool {
	return s.gc.Consume(gas, STOREBYTES, true)
}

func (s *State) GetVersioned(version int64, key StoreKey) []byte {
	_, value := s.cs.GetVersioned(version, key)
	return value
}

func (s *State) GetPrevious(num int64, key StoreKey) []byte {
	ver := s.cs.Version
	return s.GetVersioned(ver-num, key)
}
