package storage

var _ Store = &State{}
var _ Iteratable = &State{}

type State struct {
	cs        *ChainState
	cache     SessionedDirectStorage
	gc        GasCalculator
	txSession Session
}

func NewState(state *ChainState) *State {
	return &State{
		cs:    state,
		cache: NewSessionedDirectStorage(SESSION_CACHE, "state"),
		gc:    NewGasCalculator(0),
	}
}

func (s *State) WithGas(gc GasCalculator) *State {
	gs := NewGasStore(s.cache, gc)
	return &State{
		cs:    s.cs,
		cache: gs,
		gc:    gc,
	}
}

func (s *State) WithoutGas() *State {

	s.cache = NewSessionedDirectStorage(SESSION_CACHE, "state")
	s.txSession = nil
	//s.gc = NewGasCalculator(0)
	return s
}

func (s *State) BeginTxSession() {
	s.txSession = s.cache.BeginSession()
}

func (s *State) CommitTxSession() {
	if s.txSession == nil {
		panic("no tx session in state")
	}

	s.txSession.Commit()
	s.txSession = nil
}

func (s *State) DiscardTxSession() {
	s.txSession = nil
}

func (s State) Version() int64 {
	return s.cs.Version
}

func (s State) RootHash() []byte {
	return s.cs.Hash
}

func (s *State) Get(key StoreKey) ([]byte, error) {

	if s.txSession != nil {
		// Get the txSession first
		result, err := s.txSession.Get(key)
		if err == nil {
			// if got result, return directly
			return result, err
		}
	}

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
	if s.txSession != nil {
		return s.txSession.Set(key, value)
	}

	// set only for cache, waiting to be committed
	return s.cache.Set(key, value)
}

func (s *State) Exists(key StoreKey) bool {

	if s.txSession != nil {
		// check existence in txSession
		exist := s.txSession.Exists(key)
		if exist {
			return exist
		}
	}

	// check existence in cache, because it's cheaper
	exist := s.cache.Exists(key)
	if !exist {
		// if not existed in cache, check ChainState
		return s.cs.Exists(key)
	}

	return exist
}

func (s *State) Delete(key StoreKey) (bool, error) {

	if s.txSession != nil {
		return s.txSession.Delete(key)
	}
	//cache delete is always true
	_, _ = s.cache.Delete(key)

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

func (s State) Write() bool {

	s.cache.GetIterator().Iterate(func(key []byte, value []byte) bool {
		_ = s.cs.Set(key, value)
		return false
	})

	return true
}

func (s *State) Commit() (hash []byte, version int64) {

	s.Write()
	s.cache = NewSessionedDirectStorage(SESSION_CACHE, "state")
	s.txSession = nil

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
