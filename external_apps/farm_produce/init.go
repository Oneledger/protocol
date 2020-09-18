package farm_produce

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/external_apps/common"
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_action"
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_data"
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_rpc"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"os"
)

func LoadAppData(appData *common.ExtAppData) {
	logWriter := os.Stdout
	logger := log.NewLoggerWithPrefix(logWriter, "extApp").WithLevel(log.Level(4))
	//load txs
	insertProduct := common.ExtTx{
		Tx:  farm_action.InsertProduceTx{},
		Msg: &farm_action.InsertProduce{},
	}
	appData.ExtTxs = append(appData.ExtTxs, insertProduct)

	//load stores
	if dupName, ok := appData.ExtStores["extProduceStore"]; ok {
		logger.Errorf("Trying to register external store %s failed, same name already exists", dupName)
		return
	} else {
		appData.ExtStores["extProduceStore"] = farm_data.NewProduceStore(storage.NewState(appData.ChainState), "extFarmPrefix")
	}

	//load services
	balances := balance.NewStore("b", storage.NewState(appData.ChainState))
	domains := ons.NewDomainStore("ons", storage.NewState(appData.ChainState))
	olt := balance.Currency{Id: 0, Name: "OLT", Chain: chain.ONELEDGER, Decimal: 18, Unit: "nue"}
	currencies := balance.NewCurrencySet()
	err := currencies.Register(olt)
	if err != nil {
		logger.Errorf("failed to register currency %s", olt.Name, err)
		return
	}
	appData.ExtServiceMap[farm_rpc.Name()] = farm_rpc.NewService(balances, currencies, domains, logger, farm_data.NewProduceStore(storage.NewState(appData.ChainState), "extFarmPrefix"))
}
