package data

import (
	"bytes"
	"encoding/json"
)

type ContractData struct {
	Address []byte
	Data    map[string]interface{}
}

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

func (data *ContractData) GetValue(key string) interface{} {
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
func isMap(jsonMessage *json.RawMessage) bool {
	in, _ := json.Marshal(jsonMessage)
	var raw map[string]*json.RawMessage

	json.Unmarshal(in, &raw)
	return len(raw) != 0
}
