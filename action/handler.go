package action

import "github.com/tendermint/tendermint/libs/common"

type TxHandler interface {
	Validate(ctx Context, msg Msg) bool
	ProcessCheck(ctx Context, msg Msg, fee Fee) (bool, Response)
	ProcessDeliver(ctx Context, msg Msg, fee Fee) (bool, Response)
	ProcessFee(ctx Context, fee Fee) (bool, Response)
}

type Response struct {
	Data      []byte
	Log       string
	Info      string
	GasWanted int64
	GasUsed   int64
	Tags      []common.KVPair
}


