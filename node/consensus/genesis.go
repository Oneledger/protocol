package consensus

import (
	"encoding/json"
	"time"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/types"
)

type GenesisDoc = types.GenesisDoc
type GenesisValidator = types.GenesisValidator

func DefaultGenesisDoc() *GenesisDoc {
	validators := make([]GenesisValidator, 0)
	appStateBytes, err := DefaultAppState().MarshalJSON()
	if err != nil {
		log.Fatal("Failed to marshal DefaultAppState")
	}
	return &GenesisDoc{
		GenesisTime:     time.Now(),
		ChainID:         "OneLedger",
		ConsensusParams: types.DefaultConsensusParams(),
		Validators:      validators,
		AppState:        json.RawMessage(appStateBytes),
	}
}

type AppState struct {
	// Name of the account
	Account string      `json:"account"`
	States  []data.Coin `json:"states"`
}

func DefaultAppState() *AppState {
	return &AppState{
		Account: "Zero",
		States: []data.Coin{
			data.NewCoin(1000000000000, "OLT"),
			data.NewCoin(10000, "VT"),
		},
	}
}

func (a AppState) MarshalJSON() ([]byte, error) {
	states := make([]map[string]interface{}, len(a.States))
	for i := 0; i < len(a.States); i++ {
		coin := a.States[i]
		states[i] = map[string]interface{}{
			"amount":   coin.Amount.String(),
			"currency": coin.Currency.Name,
		}
	}

	jsOn := map[string]interface{}{
		"account": a.Account,
		"states":  states,
	}

	return json.Marshal(jsOn)
}

// func (a AppState) UnmarshalJSON() (()
