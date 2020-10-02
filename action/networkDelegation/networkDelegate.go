package networkDelegation

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &NetworkDelegate{}

type NetworkDelegate struct {
	UserAddress       keys.Address
	DelegationAddress keys.Address
	Amount            action.Amount
}

func (n NetworkDelegate) Signers() []action.Address {
	return []action.Address{n.UserAddress}
}

func (n NetworkDelegate) Type() action.Type {
	return action.NETWORKDELEGATE
}

func (n NetworkDelegate) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(n.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.useraddress"),
		Value: n.DelegationAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.delegationAddress"),
		Value: n.DelegationAddress.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (n NetworkDelegate) Marshal() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NetworkDelegate) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, n)
}

var _ action.Tx = networkDelegateTx{}

type networkDelegateTx struct{}

func (n networkDelegateTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	panic("implement me")
}

func (n networkDelegateTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runNetworkDelegate(ctx, tx)
}

func (n networkDelegateTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runNetworkDelegate(ctx, tx)
}

func (n networkDelegateTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.Logger.Detail("Processing Delegate Transaction for ProcessFee", signedTx)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runNetworkDelegate(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	delegate := NetworkDelegate{}
	err := delegate.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, delegate.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, delegate.Tags(), "Success")
}
