package eth

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/tendermint/tendermint/libs/common"
)

type ExtMintOETH struct {
	TrackerName  ethereum.TrackerName
	Locker action.Address
	LockAmount int64
}
var _ action.Msg = &ExtMintOETH{}

func (eem ExtMintOETH) Signers() []action.Address {
	return []action.Address{
		eem.Locker,
	}
}

func (eem ExtMintOETH) Type() action.Type {
	return action.ETH_MINT
}

func (eem ExtMintOETH) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.ETH_REPORT_FINALITY_MINT.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: eem.Locker.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(eem.TrackerName.Hex()),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (eem ExtMintOETH) Marshal() ([]byte, error) {
	return json.Marshal(eem)
}

func (eem ExtMintOETH) Unmarshal(data []byte) error {
	return json.Unmarshal(data,eem)
}


type extMintOETHTx struct {
}
func (extMintOETHTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
//Implement check Finality first

}