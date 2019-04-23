package consensus

import (
	"encoding/json"
	"time"

	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/types"
)

type GenesisDoc = types.GenesisDoc
type GenesisValidator = types.GenesisValidator

func NewGenesisDoc(chainID string) *GenesisDoc {
	validators := make([]GenesisValidator, 0)
	appStateBytes, err := NewAppState().MarshalJSON()
	if err != nil {
		log.Fatal("Failed to marshal DefaultAppState")
	}
	return &GenesisDoc{
		GenesisTime:     time.Now(),
		ChainID:         chainID,
		ConsensusParams: types.DefaultConsensusParams(),
		Validators:      validators,
		AppState:        json.RawMessage(appStateBytes),
	}
}

// TODO: Not used anymore
type State struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type AppState struct {
	// Name of the account
	Account string  `json:"account"`
	States  []State `json:"states"`
}

func NewAppState() *AppState {
	return &AppState{
		Account: "Zero",
		States: []State{
			State{Amount: "1000000000", Currency: "OLT"},
			State{Amount: "10000", Currency: "VT"},
		},
	}
}

func (a AppState) MarshalJSON() ([]byte, error) {
	states := make([]map[string]interface{}, len(a.States))
	for i := 0; i < len(a.States); i++ {
		coin := a.States[i]
		states[i] = map[string]interface{}{
			"amount":   coin.Amount,
			"currency": coin.Currency,
		}
	}

	jsOn := map[string]interface{}{
		"account": a.Account,
		"states":  states,
	}

	return json.Marshal(jsOn)
}

// func (a AppState) UnmarshalJSON() (()
