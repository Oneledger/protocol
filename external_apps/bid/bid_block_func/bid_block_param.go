package bid_block_func

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/transactions"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	abci "github.com/tendermint/tendermint/abci/types"
)

type BidParam struct {
	InternalTxStore *transactions.TransactionStore
	Logger *log.Logger
	ActionCtx action.Context
	Validator keys.Address
	Header abci.Header
	Deliver *storage.State
}
