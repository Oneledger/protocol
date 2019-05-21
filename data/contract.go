/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package data

import (
	"bytes"
	"encoding/json"
)

type ContractData struct {
	Address []byte
	Data    map[string]interface{}
}

// NewContractData returns a pointer to a new ContractData object. It takes an address as an input.
func NewContractData(address []byte) *ContractData {
	data := make(map[string]interface{})
	return &ContractData{Address: address, Data: data}
}

func (data *ContractData) Update(d map[string]interface{}) error {
	for k, v := range d {
		data.Data[k] = v
	}
	return nil
}

func (data *ContractData) Get(key string) interface{} {
	return data.Data[key]
}

func (data *ContractData) Validate(address []byte) bool {
	return bytes.Equal(address, data.Address)
}

func (data *ContractData) UpdateByJSONData(in []byte) error {
	var raw map[string]*json.RawMessage

	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	for k, v := range raw {
		value, _ := json.Marshal(v)
		data.Data[k] = value
	}
	return nil
}

//private methods
/*
func isMap(jsonMessage *json.RawMessage) bool {
	in, _ := json.Marshal(jsonMessage)
	var raw map[string]*json.RawMessage

	json.Unmarshal(in, &raw)
	return len(raw) != 0
}
*/
