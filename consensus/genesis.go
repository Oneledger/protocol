package consensus

import (
	"encoding/json"
	"github.com/Oneledger/protocol/data/fees"
	"time"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/types"
)

type GenesisDoc = types.GenesisDoc
type GenesisValidator = types.GenesisValidator

func NewGenesisDoc(chainID string, currencies balance.Currencies, feeOpt fees.FeeOption, states []StateInput) (*GenesisDoc, error) {
	validators := make([]GenesisValidator, 0)

	appStateBytes, err := newAppState(currencies, feeOpt, states).RawJSON()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal DefaultAppState")
	}
	return &GenesisDoc{
		GenesisTime:     time.Now(),
		ChainID:         chainID,
		ConsensusParams: types.DefaultConsensusParams(),
		Validators:      validators,
		AppState:        json.RawMessage(appStateBytes),
	}, nil
}

type state struct {
	// Hash of their public key
	Address string              `json:"address"`
	Balance balance.BalanceData `json:"balance"`
}

// StateInput returns a StateInput form of state, converts all serialize.Data types back into their native-form
func (s state) StateInput() StateInput {
	b := new(balance.Balance)

	// This should not return an error
	_ = b.SetData(&s.Balance)
	return StateInput{
		Address: s.Address,
		Balance: *b,
	}
}

type StateInput struct {
	Address string
	Balance balance.Balance
}

func (si StateInput) state() state {
	data := si.Balance.Data().(*balance.BalanceData)
	return state{
		Address: si.Address,
		Balance: *data,
	}
}

type AppState struct {
	Currencies balance.Currencies `json:"currencies"`
	FeeOption  fees.FeeOption     `json:"feeOption"`
	States     []state            `json:"states"`
}

func newAppState(currencies balance.Currencies, feeOpt fees.FeeOption, stateInputs []StateInput) *AppState {
	states := make([]state, len(stateInputs))
	for i, s := range stateInputs {
		states[i] = s.state()
	}

	return &AppState{
		Currencies: currencies,
		FeeOption:  feeOpt,
		States:     states,
	}
}

func (a AppState) RawJSON() ([]byte, error) {
	szr := serialize.GetSerializer(serialize.JSON)
	return szr.Serialize(a)
}
