/*
	Copyright 2017-2018 OneLedger

	Grab data from tendermint node
*/
package app

import (
	"encoding/hex"
	"github.com/Oneledger/protocol/node/global"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/privval"
)

//type PrivValidator struct {
//	Address       string    `json:"address"`
//	PubKey        TypeValue `json:"pub_key"`
//	LastHeight    int64     `json:"last_height"`
//	LastRound     int64     `json:"last_round"`
//	LastStep      int64     `json:"last_step"`
//	LastSignature string    `json:"last_signature"`
//	LastSignBytes string    `json:"last_signbytes"`
//	PrivKey       TypeValue `json:"priv_key"`
//}
//
//type TypeValue struct {
//	Type  string `json:"type"`
//	Value string `json:"value"`
//}

// Load the Priv Validator file directly from the associated Tendermint node
func LoadPrivValidatorFile() {
	filePath := global.Current.TendermintRoot + "/config/priv_validator.json"
	filepv := privval.LoadFilePV(filePath)
	address := filepv.GetAddress()
	global.Current.TendermintAddress = address.String()
	pubkey := filepv.PubKey.(ed25519.PubKeyEd25519)
	global.Current.TendermintPubKey = hex.EncodeToString(pubkey[:])

}
