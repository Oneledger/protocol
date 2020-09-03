package common

import (
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/service/query"
	"github.com/Oneledger/protocol/service/tx"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type ExtTx struct {
	Tx  action.Tx
	Msg action.Msg
}

type HandlerList []func(*ExtAppData)

var Handlers = HandlerList{}

type ExtAppData struct {
	// we need txs, data stores, services, block functions
	ExtTxs []ExtTx
	//extStores	storelist
	ExtQueryServices query.Service
	ExtTxServices    tx.Service
	//extBlockFuncs common_block.ControllerRouter
	Test string
}

func LoadExtAppData(cs *storage.ChainState) *ExtAppData {
	//this will return everything of all external apps
	//todo
	appData := &ExtAppData{}

	for _, handler := range Handlers {
		//this is to extract data from this handler to appData
		handler(appData)
	}

	return &ExtAppData{}
}

func RegisterExtApp(cs *storage.ChainState, ar action.Router, dr data.Router) error {
	//todo pass the external rpc router into here
	extAppData := LoadExtAppData(cs)
	//test
	fmt.Println(extAppData.Test)
	//register external txs using action.router
	for _, tx := range extAppData.ExtTxs {
		err := ar.AddHandler(tx.Msg.Type(), tx.Tx)
		if err != nil {
			return errors.Wrap(err, "error adding external tx")
		}
	}
	// add extStores to the data router
	//ExtAppData.extStores
	return nil
}

func (h *HandlerList) Register(fn func(*ExtAppData)) {
	*h = append(*h, fn)
}
