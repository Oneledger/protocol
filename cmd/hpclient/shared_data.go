package main

import (
	"sync"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/passport"
)

type SharedData struct {
	persons     []passport.UserID
	personAddrs []keys.Address
	mutex       sync.Mutex
}

func NewSharedData() *SharedData {
	return &SharedData{}
}

func (shd *SharedData) Lock() {
	shd.mutex.Lock()
}

func (shd *SharedData) Unlock() {
	shd.mutex.Unlock()
}
