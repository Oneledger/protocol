package data

import "github.com/Oneledger/protocol/storage"

type Store struct {
	state  *storage.State
	prefix []byte
}

type itemA struct {
	data int
}

type itemB struct {
	data string
}

type itemC struct {
	data []byte
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		state:  state,
		prefix: storage.Prefix(prefix),
	}
}

func init() {
	//Create a few stores
	db, _ := storage.GetDatabase("testDB", ".", "")

	//Create data router
	//Add stores to the data router
}
