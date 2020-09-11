package common

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/storage"
)

type ExtTx struct {
	Tx  action.Tx
	Msg action.Msg
}

type ExtServiceMap map[string]interface{}

func NewExtServiceMap() ExtServiceMap {
	return ExtServiceMap{}
}

type ExtAppData struct {
	// we need txs, data stores, services, block functions
	ChainState *storage.ChainState
	ExtTxs []ExtTx
	ExtStores	map[string]data.ExtStore
	ExtServiceMap ExtServiceMap
	ExtBlockFuncs ControllerRouter
}

func LoadExtAppData(cs *storage.ChainState) *ExtAppData {
	//this will return everything of all external apps
	appData := &ExtAppData{}
	appData.ChainState = cs
	appData.ExtStores = make(map[string]data.ExtStore)
	appData.ExtServiceMap = make(ExtServiceMap)
	functionList := make(map[txblock][]func(interface{}))
	appData.ExtBlockFuncs = ControllerRouter{
		functionlist: functionList,
	}
	for _, handler := range Handlers {
		//this is to put external app data to appData
		handler(appData)
	}

	return appData
}

func (h *HandlerList) Register(fn func(*ExtAppData)) {
	*h = append(*h, fn)
}
