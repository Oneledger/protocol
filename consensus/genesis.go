package consensus

import (
	"encoding/json"
	"time"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/network_delegation"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/types"

	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/delegation"
	ethData "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/rewards"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/serialize"
)

type GenesisValidator = types.GenesisValidator

func NewGenesisDoc(chainID string, states AppState) (*config.GenesisDoc, error) {
	validators := make([]GenesisValidator, 0)

	appStateBytes, err := states.RawJSON()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal DefaultAppState")
	}
	return &config.GenesisDoc{
		GenesisTime:     time.Now(),
		ChainID:         chainID,
		ConsensusParams: types.DefaultConsensusParams(),
		Validators:      validators,
		AppState:        json.RawMessage(appStateBytes),
		ForkParams:      config.DefaultForkParams(),
	}, nil
}

type BalanceState struct {
	Address  keys.Address   `json:"address"`
	Currency string         `json:"currency"`
	Amount   balance.Amount `json:"amount"`
}

type DomainState struct {
	Owner            keys.Address    `json:"ownerAddress"`
	Beneficiary      keys.Address    `json:"beneficiary"`
	Name             string          `json:"name"`
	CreationHeight   int64           `json:"creationHeight"`
	LastUpdateHeight int64           `json:"lastUpdateHeight"`
	ExpireHeight     int64           `json:"expireHeight"`
	ActiveFlag       bool            `json:"activeFlag"`
	OnSaleFlag       bool            `json:"onSaleFlag"`
	URI              string          `json:"uri"`
	SalePrice        *balance.Amount `json:"salePrice"`
}

//TODO: Create More Generic Struct to contain all tracker types.
type Tracker struct {
	Type          ethData.ProcessType  `json:"type"`
	State         ethData.TrackerState `json:"state"`
	TrackerName   ethchain.TrackerName `json:"trackerName"`
	SignedETHTx   []byte               `json:"signedEthTx"`
	Witnesses     []keys.Address       `json:"witnesses"`
	ProcessOwner  keys.Address         `json:"processOwner"`
	FinalityVotes []ethData.Vote       `json:"finalityVotes"`
	To            []byte               `json:"to"`
}

type ChainState struct {
	Version int64
	Hash    []byte
}

type Stake identity.Stake

type AppState struct {
	Currencies    balance.Currencies             `json:"currencies"`
	Governance    governance.GovernanceState     `json:"governance"`
	Chain         ChainState                     `json:"state"`
	Balances      []BalanceState                 `json:"balances"`
	Staking       []Stake                        `json:"staking"`
	Witness       []Stake                        `json:"witness"`
	Delegation    delegation.DelegationState     `json:"delegation"`
	Rewards       rewards.RewardMasterState      `json:"rewards"`
	Domains       []DomainState                  `json:"domains"`
	Trackers      []Tracker                      `json:"trackers"`
	Fees          []BalanceState                 `json:"fees"`
	Proposals     []governance.GovProposal       `json:"proposals"`
	NetDelegators network_delegation.State       `json:"net_delegators"`
	DelegatorRew  network_delegation.RewardState `json:"delegator_rewards"`
}

func NewAppState(currencies balance.Currencies,
	balances []BalanceState,
	staking []Stake,
	witness []Stake,
	delegation delegation.DelegationState,
	rewards rewards.RewardMasterState,
	domains []DomainState,
	fees []BalanceState,
	governance governance.GovernanceState,
) *AppState {
	return &AppState{
		Currencies: currencies,
		Balances:   balances,
		Staking:    staking,
		Witness:    witness,
		Delegation: delegation,
		Rewards:    rewards,
		Domains:    domains,
		Fees:       fees,
		Governance: governance,
	}
}

func (a AppState) RawJSON() ([]byte, error) {
	szr := serialize.GetSerializer(serialize.JSON)
	return szr.Serialize(a)
}

func GenerateState(rawState []byte) (*AppState, error) {
	state := AppState{}
	szr := serialize.GetSerializer(serialize.JSON)
	err := szr.Deserialize(rawState, state)
	if err != nil {
		return &state, err
	}
	return &state, nil
}
