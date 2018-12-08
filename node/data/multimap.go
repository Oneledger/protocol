package data

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/version"
)

type MultiMap struct {
	Name map[string]VersionMap
}

// TODO: serial won't let us use a Version struct as the map key
type VersionMap struct {
	Version map[string]Entry
}

type Entry struct {
	Value interface{}
}

func init() {
	serial.Register(MultiMap{})
	serial.Register(VersionMap{})
	serial.Register(Entry{})
}

func NewMultiMap() *MultiMap {
	return &MultiMap{
		Name: make(map[string]VersionMap, 0),
	}
}

func (mmap *MultiMap) Set(name string, version version.Version, entry interface{}) {
	var vmap VersionMap
	var ok bool

	if vmap, ok = mmap.Name[name]; !ok {
		mmap.Name[name] = *NewVersionMap()
		vmap = mmap.Name[name]
	}
	vmap.Version[version.String()] = Entry{entry}
}

func (mmap *MultiMap) Get(name string, version version.Version) Entry {
	var vmap VersionMap
	var ok bool

	if vmap, ok = mmap.Name[name]; !ok {
		return Entry{}
	}
	return vmap.Version[version.String()]
}

func NewVersionMap() *VersionMap {
	return &VersionMap{
		Version: make(map[string]Entry, 0),
	}
}
