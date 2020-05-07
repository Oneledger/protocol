package governance

import "github.com/Oneledger/protocol/storage"

type ProposalStore struct {
	state  *storage.State
	prefix []byte
}
