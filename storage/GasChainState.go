package storage

type Gas uint64

const (
	STOREBYTES Gas = 20
	READFLAT   Gas = 20
	READBYTES  Gas = 2
	WRITEFLAT  Gas = 200
	WRITEBYTES Gas = 20
	VERIFYSIG  Gas = 50
	HASHBYTES  Gas = 5
	CHECKEXIST Gas = 20
	DELETE     Gas = 50
)

// Calculate the gas used for each action, will be embeded with GasChainState.
type GasCalculator interface {
	// Consume amount of Gas for the Category
	Consume(amount, category Gas) bool

	// Get the max amount of Gas the GasCalculator accept
	GetLimit() Gas

	// Get the current consumed Gas
	GetConsumed() Gas

	// Check if the block has fullfill the Gas Limit
	IsEnough() bool
}

var _ GasCalculator = gasCalculator{}

type gasCalculator struct {
	limit    Gas
	consumed Gas
}

func (g gasCalculator) IsEnough() bool {
	panic("implement me")
}

func (g gasCalculator) Consume(amount, category Gas) bool {
	panic("implement me")
}

func (g gasCalculator) GetLimit() Gas {
	panic("implement me")
}

func (g gasCalculator) GetConsumed() Gas {
	panic("implement me")
}

func NewGasCalculator(limit Gas) GasCalculator {
	return &gasCalculator{
		limit:    limit,
		consumed: 0,
	}
}

type GasChainState struct {
	*ChainState
	GasCalculator
}
