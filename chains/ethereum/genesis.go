package eth

type AppState struct {
	// The contract address for the smart contract
	ContractAddress Address `json:"contractAddress"`
	// The initial validator set for the contract
	InitialValidators []Address `json:"initialValidators"`
	// Is Ethereum enabled on this chain? If so, validators will fail initialization
	// if they lack of a valid ethereum key
}

func DefaultAppState() *AppState {
	return new(AppState)
}

func NewDisabledAppState() *AppState {
	return nil
}

// NewAppState returns a new enabled AppState with the given parameters
// this should be called on init
func NewAppState(contract Address, initialValidators []Address) AppState {
	return AppState{
		ContractAddress:   contract,
		InitialValidators: initialValidators,
	}
}

// Contract returns a contract object for interaction
func (state *AppState) Contract() (*Contract, error) {
	return nil, ErrNotImplemented
}
