package common

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/service/query"
	"github.com/Oneledger/protocol/service/tx"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type ExtTx struct {
	Tx action.Tx
	Msg action.Msg
}

var sfddsfs handler list

type ExtAppData struct {
	// we need txs, data stores, services, block functions
	extTxs	[]ExtTx
	extStores	storelist
	extQueryServices query.Service
	extTxServices tx.Service
	//extBlockFuncs common_block.ControllerRouter
}

func LoadExtAppData(cs *storage.ChainState) *ExtAppData{
	//todo
	return &ExtAppData{}
}

func RegisterExtApp(cs *storage.ChainState, ar action.Router, dr data.Router) error {
	ExtAppData := LoadExtAppData(cs)
	//register external txs using action.router
	for _, tx := range ExtAppData.extTxs {
		err := ar.AddHandler(tx.Msg.Type(), tx.Tx)
		if err != nil {
			return errors.Wrap(err, "error adding external tx")
		}
	}
	// add extStores to the data router
	ExtAppData.extStores
	return nil
}