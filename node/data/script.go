package data

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/version"
)

type ScriptRecords struct {
	Name    string          //Preserved, not used
	Version version.Version //PReserved, not used
	Script  Script
}

type ContractRefRecord struct {
	Status ContractRefStatus
}

type ContractRefStatus = uint32 //smart contract ref status

const (
	PENDING   ContractRefStatus = 100
	COMPLETED ContractRefStatus = 200
	NOT_FOUND ContractRefStatus = 404
	ERROR     ContractRefStatus = 500
)

type Script struct { ///code script for smart contract
	Script []byte
}

func init() {
	serial.Register(ScriptRecords{})
	serial.Register(Script{})
	serial.Register(ContractRefRecord{})
}

func NewScriptRecords() *ScriptRecords {
	return &ScriptRecords{}
}

func (scriptRecords *ScriptRecords) Set(name string, version version.Version, script Script) {
	scriptRecords.Name = name
	scriptRecords.Version = version
	scriptRecords.Script = script
}
