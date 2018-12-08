package data

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/version"
)

type ScriptRecords struct {
	Name map[string]Versions
}

type Versions struct {
	Version map[string]Script
}

type Script struct {
	Script []byte
}

func init() {
	serial.Register(ScriptRecords{})
	serial.Register(Versions{})
	serial.Register(Script{})
}

func NewScriptRecords() *ScriptRecords {
	return &ScriptRecords{
		Name: make(map[string]Versions, 0),
	}
}

func (scriptRecords *ScriptRecords) Set(name string, version version.Version, script Script) {
	var versions Versions
	var ok bool

	if versions, ok = scriptRecords.Name[name]; !ok {
		scriptRecords.Name[name] = *NewVersions()
		versions = scriptRecords.Name[name]
	}
	versions.Version[version.String()] = script
}

func NewVersions() *Versions {
	return &Versions{
		Version: make(map[string]Script, 0),
	}
}
