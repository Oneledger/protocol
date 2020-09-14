package external_apps

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/external_apps/common"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

func init() {
	//register new external app handler function in the last line
}

func RegisterExtApp(cs *storage.ChainState, ar action.Router, dr data.Router, esm common.ExtServiceMap, cr common.ControllerRouter) error {
	extAppData := common.LoadExtAppData(cs)
	//register external txs using action.router
	for _, tx := range extAppData.ExtTxs {
		err := ar.AddHandler(tx.Msg.Type(), tx.Tx)
		if err != nil {
			return errors.Wrap(err, "error adding external tx")
		}
	}

	// add extStores to app data
	for name, store := range extAppData.ExtStores {
		err := dr.Add(data.Type(name), store)
		if err != nil {
			return errors.Wrap(err, "error adding external store")
		}
	}
	// add services
	for name, service := range extAppData.ExtServiceMap {
		esm[name] = service
	}

	//add block beginner & ender function router here
	beginnerFuncs, _ := extAppData.ExtBlockFuncs.Iterate(common.BlockBeginner)
	for _, function := range beginnerFuncs {
		err := cr.Add(common.BlockBeginner, function)
		if err != nil {
			return errors.Wrap(err, "error adding external block beginner funcs")
		}
	}
	enderFuncs, _ := extAppData.ExtBlockFuncs.Iterate(common.BlockEnder)
	for _, function := range enderFuncs {
		err := cr.Add(common.BlockEnder, function)
		if err != nil {
			return errors.Wrap(err, "error adding external block ender funcs")
		}
	}
	return nil
}

