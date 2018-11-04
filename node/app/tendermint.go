/*
	Copyright 2017-2018 OneLedger

	Grab data from tendermint node
*/
package app

import (
	"encoding/json"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"io/ioutil"
	"os"
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

func LoadPrivValidatorFile() {

	log.Debug("LoadPrivValidatorFile", "global.Current.TendermintRoot", global.Current.TendermintRoot)
	filePath := global.Current.TendermintRoot + "/config/priv_validator.json"
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Debug("FeePayment", "err", err)
	}
	log.Debug("FeePaymentDat", "jsonFile", jsonFile)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var privValidator PrivValidator

	json.Unmarshal(byteValue, &privValidator)
	log.Debug("FeePaymentValAddress", "address", privValidator.Address)
	log.Debug("FeePaymentValPubKey", "ValPubKey", privValidator.PubKey.Value)
	global.Current.TendermintAddress = privValidator.Address
	global.Current.TendermintPubKey = privValidator.PubKey.Value

}
