package bid

import (
	"fmt"
	"github.com/Oneledger/protocol/external_apps/common"
)

type AppData struct {
	// we need txs, data stores, services, block functions
	//ExtTxs	[]common.ExtTx
	//extStores	storelist
	//ExtQueryServices query.Service
	//ExtTxServices tx.Service
	//extBlockFuncs common_block.ControllerRouter
	Test string
}

func init() {
	common.Handlers.Register(LoadAppData)
}

func LoadAppData(appData *common.ExtAppData) {

	appData.Test = "TEST BID APPLICATION"
	fmt.Println("LOADING BID DATA !!!")
	//this will include everything about this external app
	//extxs := make([]common.ExtTx, 10)
	//bidCreate := common.ExtTx{
	//	Tx: bid_action.CreateBidTx{},
	//	Msg: &bid_action.CreateBid{},
	//}
	//append(extxs, bid_action.BID_CREATE)

	//load txs
	//load services
	//load stores
	//load beginner and ender functions
}
