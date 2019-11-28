package eth

import (
	"encoding/json"

	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
)

var _ action.Msg = &Redeem{}

type Redeem struct {
	Owner  action.Address
	To     action.Address
	ETHTxn []byte
}

func (r Redeem) Signers() []action.Address {
	return []action.Address{r.Owner}
}

func (r Redeem) Type() action.Type {
	return action.ETH_REDEEM
}

func (r Redeem) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(r.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: r.Owner,
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (r Redeem) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Redeem) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

var _ action.Tx = ethRedeemTx{}

type ethRedeemTx struct {
}

func (ethRedeemTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	panic("implement me")
}

func (ethRedeemTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	panic("implement me")
}

func (ethRedeemTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	panic("implement me")
}

func (ethRedeemTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	panic("implement me")
}
