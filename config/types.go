package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

type GenesisDoc struct {
	GenesisTime     time.Time                  `json:"genesis_time"`
	ChainID         string                     `json:"chain_id"`
	ConsensusParams *tmtypes.ConsensusParams   `json:"consensus_params,omitempty"`
	Validators      []tmtypes.GenesisValidator `json:"validators,omitempty"`
	AppHash         tmbytes.HexBytes           `json:"app_hash"`
	AppState        json.RawMessage            `json:"app_state,omitempty"`
	ForkParams      *ForkParams                `json:"fork"`
}

// SaveAs is a utility method for saving GenensisDoc as a JSON file.
func (genDoc *GenesisDoc) SaveAs(file string) error {
	genDocBytes, err := tmtypes.GetCodec().MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}
	return tmos.WriteFile(file, genDocBytes, 0644)
}

// ValidatorHash returns the hash of the validator set contained in the GenesisDoc
func (genDoc *GenesisDoc) ValidatorHash() []byte {
	vals := make([]*tmtypes.Validator, len(genDoc.Validators))
	for i, v := range genDoc.Validators {
		vals[i] = tmtypes.NewValidator(v.PubKey, v.Power)
	}
	vset := tmtypes.NewValidatorSet(vals)
	return vset.Hash()
}

// ValidateAndComplete checks that all necessary fields are present
// and fills in defaults for optional fields left empty
func (genDoc *GenesisDoc) ValidateAndComplete() error {
	if genDoc.ChainID == "" {
		return errors.New("genesis doc must include non-empty chain_id")
	}
	if len(genDoc.ChainID) > tmtypes.MaxChainIDLen {
		return errors.Errorf("chain_id in genesis doc is too long (max: %d)", tmtypes.MaxChainIDLen)
	}

	if genDoc.ConsensusParams == nil {
		genDoc.ConsensusParams = tmtypes.DefaultConsensusParams()
	} else if err := genDoc.ConsensusParams.Validate(); err != nil {
		return err
	}

	if genDoc.ForkParams == nil {
		genDoc.ForkParams = DefaultForkParams()
	} else if err := genDoc.ForkParams.Validate(); err != nil {
		return err
	}

	for i, v := range genDoc.Validators {
		if v.Power == 0 {
			return errors.Errorf("the genesis file cannot contain validators with no voting power: %v", v)
		}
		if len(v.Address) > 0 && !bytes.Equal(v.PubKey.Address(), v.Address) {
			return errors.Errorf("incorrect address for validator %v in the genesis file, should be %v", v, v.PubKey.Address())
		}
		if len(v.Address) == 0 {
			genDoc.Validators[i].Address = v.PubKey.Address()
		}
	}

	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = tmtime.Now()
	}

	return nil
}

//------------------------------------------------------------
// Make genesis state from file

// GenesisDocFromJSON unmarshalls JSON data into a GenesisDoc.
func GenesisDocFromJSON(jsonBlob []byte) (*GenesisDoc, error) {
	genDoc := GenesisDoc{}
	err := tmtypes.GetCodec().UnmarshalJSON(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	return &genDoc, err
}

// GenesisDocFromFile reads JSON data from a file and unmarshalls it into a GenesisDoc.
func GenesisDocFromFile(genDocFile string) (*GenesisDoc, error) {
	jsonBlob, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read GenesisDoc file")
	}
	genDoc, err := GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error reading GenesisDoc at %v", genDocFile))
	}
	return genDoc, nil
}

// ForkParams determine the fork blocks number where to apply the global update for network
type ForkParams struct {
	FrankensteinBlock int64 `json:"frankensteinBlock"`
}

// DefaultForkParams initial config
func DefaultForkParams() *ForkParams {
	return &ForkParams{
		FrankensteinBlock: 1, // 0 means disabled as tendermint blocks started from 1
	}
}

// ToMap converts fork struct to map
func (f *ForkParams) ToMap() (map[string]interface{}, error) {
	var inMap map[string]interface{}
	inrec, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(inrec, &inMap)
	return inMap, err
}

// IsFrankensteinBlock check if fork update arrived to apply change at specific block
func (f *ForkParams) IsFrankensteinBlock(height int64) bool {
	return f.FrankensteinBlock != 0 && f.FrankensteinBlock == height
}

// IsFrankensteinUpdate check if fork update arrived to apply changes after specific block
func (f *ForkParams) IsFrankensteinUpdate(height int64) bool {
	return f.FrankensteinBlock != 0 && f.FrankensteinBlock <= height
}

// Validate validates the ForkParams to ensure all values are within their
// allowed limits, and returns an error if they are not.
func (f *ForkParams) Validate() error {
	return nil
}
