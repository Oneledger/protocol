package data

import "github.com/Oneledger/protocol/storage"

type State struct {
	ChainState *storage.ChainState
	Stores     map[string]storage.Store
}

func (s State) Commit() ([]byte, int64) {
	return s.ChainState.Commit()
}
