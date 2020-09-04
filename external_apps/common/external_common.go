package common

import (
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/service/query"
	"github.com/Oneledger/protocol/service/tx"
	"github.com/Oneledger/protocol/storage"
)

type ExtTx struct {
	Tx  action.Tx
	Msg action.Msg
}

type ExtAppData struct {
	// we need txs, data stores, services, block functions
	ChainState *storage.ChainState
	ExtTxs []ExtTx
	ExtStores	map[string]interface{}
	ExtServiceMap map[string]interface{}
	ExtQueryServices query.Service
	ExtTxServices    tx.Service
	//extBlockFuncs common_block.ControllerRouter
	Test string
}

func LoadExtAppData(cs *storage.ChainState) *ExtAppData {
	//this will return everything of all external apps
	//todo
	appData := &ExtAppData{}
	appData.ChainState = cs
	fmt.Println("LoadExtAppData runs, and the Handlers: ", Handlers)
	for _, handler := range Handlers {
		//this is to extract data from this handler to appData
		fmt.Println("inside for loop")
		handler(appData)
	}

	return appData
}

func (h *HandlerList) Register(fn func(*ExtAppData)) {
	*h = append(*h, fn)
}
