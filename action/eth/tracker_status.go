package eth

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/tendermint/tendermint/libs/common"
)

type trackerquery struct {
	trackername ethereum.TrackerName
}

func (t trackerquery) Signers() []action.Address {
	panic("implement me")
}

func (t trackerquery) Type() action.Type {
	panic("implement me")
}

func (t trackerquery) Tags() common.KVPairs {
	panic("implement me")
}

func (t trackerquery) Marshal() ([]byte, error) {
	panic("implement me")
}

func (t trackerquery) Unmarshal([]byte) error {
	panic("implement me")
}
