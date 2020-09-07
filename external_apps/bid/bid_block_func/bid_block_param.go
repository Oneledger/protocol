package bid_block_func

import (
	"github.com/Oneledger/protocol/data/transactions"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
	"github.com/Oneledger/protocol/log"
	abci "github.com/tendermint/tendermint/abci/types"
)

type BidParam struct {
	BidMasterStore *bid_data.BidMasterStore
	InternalTxStore *transactions.TransactionStore
	Logger *log.Logger
	Header abci.Header

}
