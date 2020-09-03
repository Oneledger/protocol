package external_apps

import (
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/external_apps/bid"
	"github.com/Oneledger/protocol/external_apps/common"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

func init() {
	fmt.Println("init from externalApps/init")
	common.Handlers.Register(bid.LoadAppData)
}

func RegisterExtApp(cs *storage.ChainState, ar action.Router, dr data.Router) error {
	//todo pass the external rpc router into here
	extAppData := common.LoadExtAppData(cs)
	//test
	fmt.Println("extAppData.Test", extAppData.Test)
	//register external txs using action.router
	for _, tx := range extAppData.ExtTxs {
		err := ar.AddHandler(tx.Msg.Type(), tx.Tx)
		if err != nil {
			return errors.Wrap(err, "error adding external tx")
		}
	}

	// add extStores to the data router
	//for _, stores := range extAppData.
	//ExtAppData.extStores
	return nil
}

