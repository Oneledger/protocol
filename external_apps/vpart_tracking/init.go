package vpart_tracking

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/external_apps/common"
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_action"
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_data"
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_rpc"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"os"
)

func LoadAppData(appData *common.ExtAppData) {
	logWriter := os.Stdout
	logger := log.NewLoggerWithPrefix(logWriter, "extApp").WithLevel(log.Level(4))
	//load txs
	insertProduce := common.ExtTx{
		Tx:  vpart_action.InsertPartTx{},
		Msg: &vpart_action.InsertPart{},
	}
	appData.ExtTxs = append(appData.ExtTxs, insertProduce)

	//load stores
	if dupName, ok := appData.ExtStores["extVPartStore"]; ok {
		logger.Errorf("Trying to register external store %s failed, same name already exists", dupName)
		return
	} else {
		appData.ExtStores["extVPartStore"] = vpart_data.NewVPartStore(storage.NewState(appData.ChainState), "extVPartPrefix")
	}

	//load services
	balances := balance.NewStore("b", storage.NewState(appData.ChainState))
	olt := balance.Currency{Id: 0, Name: "OLT", Chain: chain.ONELEDGER, Decimal: 18, Unit: "nue"}
	currencies := balance.NewCurrencySet()
	err := currencies.Register(olt)
	if err != nil {
		logger.Errorf("failed to register currency %s", olt.Name, err)
		return
	}
	appData.ExtServiceMap[vpart_rpc.Name()] = vpart_rpc.NewService(balances, currencies, logger, vpart_data.NewVPartStore(storage.NewState(appData.ChainState), "extVPartPrefix"))
}
