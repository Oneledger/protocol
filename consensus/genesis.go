package consensus

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/types"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/serialize"
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

type GovernanceState struct {
	FeeOption   fees.FeeOption             `json:"feeOption"`
	ETHCDOption ethchain.ChainDriverOption `json:"ethchaindriverOption"`
	BTCCDOption bitcoin.ChainDriverOption  `json:"bitcoinChainDriverOption"`
	ONSOptions  ons.Options                `json:"onsOptions"`
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

type ChainState struct {
	Version int64
	Hash    []byte
}

type Stake identity.Stake

type AppState struct {
	Currencies balance.Currencies `json:"currencies"`
	Governance GovernanceState    `json:"governance"`
	Chain      ChainState         `json:"state"`
	Balances   []BalanceState     `json:"balances"`
	Staking    []Stake            `json:"staking"`
	Domains    []DomainState      `json:"domains"`
	Fees       []BalanceState     `json:"fees"`
}

func NewAppState(currencies balance.Currencies,
	balances []BalanceState,
	staking []Stake,
	domains []DomainState,
	fees []BalanceState,
	governance GovernanceState,
) *AppState {
	return &AppState{
		Currencies: currencies,
		Balances:   balances,
		Staking:    staking,
		Domains:    domains,
		Fees:       fees,
		Governance: governance,
	}
}

func (a AppState) RawJSON() ([]byte, error) {
	szr := serialize.GetSerializer(serialize.JSON)
	return szr.Serialize(a)
}
