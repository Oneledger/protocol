package bid

import (
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
	common.HandlerList.register(loadAppData)
}

func loadAppData() AppData {
	//this will include everything about this external app
	//extxs := make([]common.ExtTx, 10)
	//bidCreate := common.ExtTx{
	//	Tx: bid_action.CreateBidTx{},
	//	Msg: &bid_action.CreateBid{},
	//}
	//append(extxs, bid_action.BID_CREATE)
	return AppData{
		Test: "my test",
	}
	//load txs
	//load services
	//load stores
	//load beginner and ender functions
}
