package storage

type Gas int

const (
	STOREBYTES Gas = 20
	READFLAT   Gas = 20
	READBYTES  Gas = 2
	WRITEFLAT  Gas = 200
	WRITEBYTES Gas = 20
	VERIFYSIG  Gas = 5000
	HASHBYTES  Gas = 5
	CHECKEXIST Gas = 20
	DELETE     Gas = 50
	FLAT       Gas = 1
	CONTRACT   Gas = 1
)

// Calculate the gas used for each action, will be embedded with GasStore.
type GasCalculator interface {
	// Consume amount of Gas for the Category
	Consume(amount, category Gas, allowOverflow bool) bool

	// Get the max amount of Gas the GasCalculator accept
	GetLimit() Gas

	// Get the current consumed Gas
	GetConsumed() Gas

	// Check if the block has fullfill the Gas Limit
	IsEnough() bool

	// GetLeft return total gas left, used right now in olvm (frankenstein update)
	GetLeft() uint64
}

var _ GasCalculator = &gasCalculator{}

type gasCalculator struct {
	limit    Gas
	consumed Gas
}

func (g gasCalculator) IsEnough() bool {
	if g.consumed >= g.limit {
		return true
	}
	return false
}

func (g *gasCalculator) Consume(amount, category Gas, allowOverflow bool) bool {
	currentGasCost := amount * category
	if allowOverflow {
		g.consumed += currentGasCost
		return true
	} else {
		if g.consumed >= g.limit {
			return false
		} else {
			g.consumed += currentGasCost
			return true
		}
	}
}

func (g gasCalculator) GetLimit() Gas {
	return g.limit
}

func (g gasCalculator) GetConsumed() Gas {
	return g.consumed
}

func (g gasCalculator) GetLeft() uint64 {
	var availableGas Gas

	if g.consumed >= g.limit {
		availableGas = 0
	} else {
		availableGas = g.limit - g.consumed
	}
	return uint64(availableGas)
}

func NewGasCalculator(limit Gas) GasCalculator {
	return &gasCalculator{
		limit:    limit,
		consumed: 0,
	}
}

type IGasStore interface {
	Set(key StoreKey, value []byte) error
	Get(key StoreKey) ([]byte, error)
	Exists(key StoreKey) bool
	Delete(key StoreKey) (bool, error)

	BeginSession() Session
	Close()
	DumpState()
	GetIterable() Iterable
}

type GasStore struct {
	SessionedDirectStorage
	GasCalculator
}

func NewGasStore(store SessionedDirectStorage, gc GasCalculator) *GasStore {
	return &GasStore{
		SessionedDirectStorage: store,
		GasCalculator:          gc,
	}
}

func (g *GasStore) Set(key StoreKey, value []byte) error {
	ok := g.GasCalculator.Consume(Gas(1), WRITEFLAT, false)
	if !ok {
		return ErrExceedGasLimit
	}
	err := g.SessionedDirectStorage.Set(key, value)
	if err != nil {
		return err
	}
	g.GasCalculator.Consume(Gas(len(value)), WRITEBYTES, true)
	return nil
}

func (g *GasStore) Get(key StoreKey) ([]byte, error) {
	ok := g.GasCalculator.Consume(Gas(1), READFLAT, false)
	if !ok {
		//log.Error(ErrExceedGasLimit.Error())
		return nil, ErrExceedGasLimit
	}
	value, err := g.SessionedDirectStorage.Get(key)
	if err != nil {
		return nil, err
	}
	ok = g.GasCalculator.Consume(Gas(len(value)), READBYTES, true)
	return value, nil
}

func (g *GasStore) Exists(key StoreKey) bool {
	ok := g.GasCalculator.Consume(Gas(1), CHECKEXIST, false)
	if !ok {
		log.Error(ErrExceedGasLimit.Error())
		return false
	}

	exist := g.SessionedDirectStorage.Exists(key)
	return exist
}

func (g *GasStore) Delete(key StoreKey) (bool, error) {
	ok := g.GasCalculator.Consume(Gas(1), DELETE, false)
	if !ok {
		log.Error(ErrExceedGasLimit.Error())
		return false, ErrExceedGasLimit
	}

	return g.SessionedDirectStorage.Delete(key)
}

// NoGasStore skip gas consume because EVM will calculate it and apply at fee processor
type NoGasStore struct {
	SessionedDirectStorage
	GasCalculator
}

func NewNoGasStore(store SessionedDirectStorage) *NoGasStore {
	return &NoGasStore{
		SessionedDirectStorage: store,
		GasCalculator:          nil,
	}
}

func (g *NoGasStore) Set(key StoreKey, value []byte) error {
	err := g.SessionedDirectStorage.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (g *NoGasStore) Get(key StoreKey) ([]byte, error) {
	value, err := g.SessionedDirectStorage.Get(key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (g *NoGasStore) Exists(key StoreKey) bool {
	exist := g.SessionedDirectStorage.Exists(key)
	return exist
}

func (g *NoGasStore) Delete(key StoreKey) (bool, error) {
	return g.SessionedDirectStorage.Delete(key)
}
