package action

import "github.com/tendermint/tendermint/libs/common"

type Tx interface {
	//it should be able to validate a tx by itself, non-valid tx will be reject by the node directly
	//without going to node process
	Validate(msg Msg, fee Fee, signatures []Signature) (bool, error)

	//check tx on the first node who receives it, if process check failed, the tx will not be broadcast
	//could store a version of checked value in storage, but not implemented now
	ProcessCheck(ctx *Context, msg Msg, fee Fee) (bool, Response)

	//deliver tx on the chain by changing the storage values, which should only be committed at Application
	//commit stage
	ProcessDeliver(ctx *Context, msg Msg, fee Fee) (bool, Response)

	//process the charge of fees
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
