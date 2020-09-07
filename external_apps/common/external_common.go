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

type ExtServiceMap map[string]interface{}

func NewExtServiceMap() *ExtServiceMap {
	return &ExtServiceMap{}
}

type ExtAppData struct {
	// we need txs, data stores, services, block functions
	ChainState *storage.ChainState
	ExtTxs []ExtTx
	ExtStores	map[string]interface{}
	ExtServiceMap ExtServiceMap
	ExtQueryServices query.Service
	ExtTxServices    tx.Service
	ExtBlockFuncs ControllerRouter
	Test string
}

func LoadExtAppData(cs *storage.ChainState) *ExtAppData {
	//this will return everything of all external apps
	appData := &ExtAppData{}
	appData.ChainState = cs
	fmt.Println("LoadExtAppData runs, and the Handlers: ", Handlers)
	for _, handler := range Handlers {
		//this is to put external app data to appData
		fmt.Println("inside for loop")
		handler(appData)
	}

	return appData
}

func (h *HandlerList) Register(fn func(*ExtAppData)) {
	*h = append(*h, fn)
}
