/*
	Copyright 2017-2018 OneLedger

	Declare basic types used by the Application

	If a type requires functions or a few types are intertwinded, then should be in their own file.
*/
package action

import (
	"bytes"

	"strconv"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

func init() {
	serial.Register(Event{})
}

func GetHeight(app interface{}) int64 {
	balances := GetBalances(app)

	height := int64(balances.Version)
	return height
}

type Event struct {
	Type        Type   `json:"type"`
	SwapKeyHash []byte `json:"swapkeyhash"`
	Step        int    `json:"step"`
}

func (e Event) ToKey() []byte {
	buffer, err := serial.Serialize(e, serial.CLIENT)
	if err != nil {
		log.Error("Failed to Serialize event key")
	}
	return buffer
}

func SaveEvent(app interface{}, eventKey Event, status bool) {
	events := GetEvent(app)
	s := "0"
	if status {
		s = "1"
	}
	log.Debug("Save Event", "key", eventKey)

	session := events.Begin()
	session.Set(eventKey.ToKey(), []byte(s))
	session.Commit()
}

func FindEvent(app interface{}, eventKey Event) bool {
	log.Debug("Load Event", "key", eventKey)
	events := GetEvent(app)
	result := events.Get(eventKey.ToKey())
	if result == nil {
		return false
	}

	if bytes.Equal(result.([]byte), []byte("1")) {
		return true
	}

	return false
}

func SaveContract(app interface{}, contractKey []byte, nonce int64, contract []byte) {
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	log.Debug("Save contract", "key", contractKey, "afterNonce", n)
	session := contracts.Begin()
	session.Set(n, contract)
	session.Commit()
}

func FindContract(app interface{}, contractKey []byte, nonce int64) []byte {
	log.Debug("Load Contract", "key", contractKey)
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	result := contracts.Get(n)
	if result == nil {
		return nil
	}
	return result.([]byte)
}

func DeleteContract(app interface{}, contractKey []byte, nonce int64) {
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	log.Debug("Delete contract", "key", contractKey, "afterNonce", n)
	session := contracts.Begin()
	session.Delete(n)
	session.Commit()
}
