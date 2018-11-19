/*
	Copyright 2017-2018 OneLedger

	Grab data from tendermint node
*/
package app

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
)

type PrivValidator struct {
	Address       string    `json:"address"`
	PubKey        TypeValue `json:"pub_key"`
	LastHeight    int64     `json:"last_height"`
	LastRound     int64     `json:"last_round"`
	LastStep      int64     `json:"last_step"`
	LastSignature string    `json:"last_signature"`
	LastSignBytes string    `json:"last_signbytes"`
	PrivKey       TypeValue `json:"priv_key"`
}

type TypeValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Load the Priv Validator file directly from the associated Tendermint node
func LoadPrivValidatorFile() {
	filePath := global.Current.TendermintRoot + "/config/priv_validator.json"
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Debug("FeePayment", "earr", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var privValidator PrivValidator

	json.Unmarshal(byteValue, &privValidator)
	global.Current.TendermintAddress = privValidator.Address
	global.Current.TendermintPubKey = privValidator.PubKey.Value

}
