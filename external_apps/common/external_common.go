package common

import (
	"github.com/Oneledger/protocol/external_apps/common/common_action"
	"github.com/Oneledger/protocol/external_apps/common/common_data"
	"github.com/Oneledger/protocol/external_apps/common/common_rpc/common_rpc_query"
	"github.com/Oneledger/protocol/storage"
)

type ExtAppData struct {
	// we need txs, data stores, services, block functions
	extTxs	[]common_action.Msg
	extStores	common_data.Router
	extQueryServices common_rpc_query.Service
	extTxServices common_rpc_query.Service
	extBlockBeginnerFuncs
	extBlockEnderFuncs
}

func LoadExtAppData(cs *storage.ChainState) *ExtAppData{

	return &ExtAppData{}
}