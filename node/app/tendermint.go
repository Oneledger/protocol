/*
	Copyright 2017-2018 OneLedger

	Grab data from tendermint node
*/
package app

import (
	"encoding/hex"
	"path/filepath"

	"github.com/Oneledger/protocol/node/global"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/privval"
)

// Load the Priv Validator file directly from the associated Tendermint node
func LoadPrivValidatorFile() {
	keyfilePath := filepath.Join(global.Current.ConsensusDir(), "config", "priv_validator_key.json")
	statefilePath := filepath.Join(global.Current.ConsensusDir(), "data", "priv_validator_state.json")
	filepv := privval.LoadFilePV(keyfilePath, statefilePath)
	address := filepv.GetAddress()
	global.Current.TendermintAddress = address.String()
	pubkey := filepv.GetPubKey().(ed25519.PubKeyEd25519)
	global.Current.TendermintPubKey = hex.EncodeToString(pubkey[:])

}
