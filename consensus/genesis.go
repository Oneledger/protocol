package consensus

import (
	"encoding/json"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"time"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/types"
)

type GenesisDoc = types.GenesisDoc
type GenesisValidator = types.GenesisValidator

func NewGenesisDoc(chainID string, states AppState) (*GenesisDoc, error) {
	validators := make([]GenesisValidator, 0)

	appStateBytes, err := states.RawJSON()
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

type BalanceState struct {
	Address  keys.Address   `json:"address"`
	Currency string         `json:"currency"`
	Amount   balance.Amount `json:"amount"`
}

type DomainState struct {
	OwnerAddress   keys.Address `json:"ownerAddress"`
	AccountAddress keys.Address `json:"accountAddress"`
	Name           string       `json:"name"`
}

type Stake identity.Stake

type AppState struct {
	Currencies balance.Currencies `json:"currencies"`
	FeeOption  fees.FeeOption     `json:"feeOption"`
	Balances   []BalanceState     `json:"balances"`
	Staking    []Stake            `json:"staking"`
	Domains    []DomainState      `json:"domains"`
	Fees       []BalanceState     `json:"fees"`
}

func NewAppState(currencies balance.Currencies,
	feeOpt fees.FeeOption,
	balances []BalanceState,
	staking []Stake,
	domains []DomainState,
	fees []BalanceState,
) *AppState {
	return &AppState{
		Currencies: currencies,
		FeeOption:  feeOpt,
		Balances:   balances,
		Staking:    staking,
		Domains:    domains,
		Fees:       fees,
	}
}

func (a AppState) RawJSON() ([]byte, error) {
	szr := serialize.GetSerializer(serialize.JSON)
	return szr.Serialize(a)
}
