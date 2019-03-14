package data

import (
	"testing"
)

func TestCreateNewContractData(t *testing.T) {
	address := "0x0"
	_ = NewContractData([]byte(address))
}

func TestUpdate(t *testing.T) {
	address := "0x0"
	contractData := NewContractData([]byte(address))
	var toUpdate = make(map[string]interface{})
	toUpdate["key1"] = "value1"
	toUpdate["key2"] = 12
	contractData.Update(toUpdate)
	if contractData.Data["key1"] != "value1" {
		t.Error("update key 1 failed")
	}

	if contractData.Data["key2"] != 12 {
		t.Error("update key 2 failed")
	}
}

func TestPartialUpdate(t *testing.T) {
	contractData := NewContractData([]byte("0x0"))
	contractData.Data["key1"] = "originalValue"
	contractData.Data["key3"] = "value3"
	var toUpdate = make(map[string]interface{})
	toUpdate["key1"] = "value1"
	toUpdate["key2"] = 12
	contractData.Update(toUpdate)
	if contractData.Data["key1"] != "value1" {
		t.Error("update key 1 failed")
	}

	if contractData.Data["key2"] != 12 {
		t.Error("update key 2 failed")
	}

	contractData.Update(toUpdate)
	if contractData.Data["key1"] != "value1" {
		t.Error("update key 1 failed")
	}

	if contractData.Data["key3"] != "value3" {
		t.Error("update key3 by mistake")
	}
}

func TestGet(t *testing.T) {
	contractData := NewContractData([]byte("0x0"))
	contractData.Data["key1"] = "value1"
	contractData.Data["key2"] = "value2"
	if contractData.Get("key1") != "value1" {
		t.Error("cannot get the value by the key")
	}

	if contractData.Get("key3") != nil {
		t.Error("cannot get the null with invalid key")
	}
}

func TestValidate(t *testing.T) {
	contractData := NewContractData([]byte("0x0"))
	if !contractData.Validate([]byte("0x0")) {
		t.Error("cannot validate the contract data")
	}
}
func TestUpdateByJSONData(t *testing.T) {
	contractData := NewContractData([]byte("0x0"))
	contractData.Data["key1"] = "value1"
	contractData.Data["key2"] = "value2"
	jsonStr := `{"key1": "new value 1", "key3" : {"key3_1" : "key3_2"}}`
	contractData.UpdateByJSONData([]byte(jsonStr))
	if string(contractData.Get("key1").([]byte)) != `"new value 1"` {
		t.Error("update key1 failed")
	}
	if contractData.Get("key2") != "value2" {
		t.Error("update key2 by mistake ")
	}

	if string(contractData.Get("key3").([]byte)) != `{"key3_1":"key3_2"}` {
		t.Error("update key3 failed ")
	}
}
