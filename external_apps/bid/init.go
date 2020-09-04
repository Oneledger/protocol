package bid

import (
	"fmt"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/external_apps/bid/bid_action"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
	"github.com/Oneledger/protocol/external_apps/bid/bid_rpc/bid_rpc_query"
	"github.com/Oneledger/protocol/external_apps/common"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/service/query"
	"github.com/Oneledger/protocol/storage"
)

var logger *log.Logger

func init() {
	//fmt.Println("init from bid/init")
	common.Handlers.Register(LoadAppData)
}

//this is the handler function, it will add a bunch of things into appData
func LoadAppData(appData *common.ExtAppData) {

	appData.Test = "TEST BID APPLICATION"
	fmt.Println("LOADING BID DATA !!!")
	//load txs

	bidCreate := common.ExtTx{
		Tx: bid_action.CreateBidTx{},
		Msg: &bid_action.CreateBid{},
	}
	bidCancel := common.ExtTx{
		Tx: bid_action.CancelBidTx{},
		Msg: &bid_action.CancelBid{},
	}
	bidExpire := common.ExtTx{
		Tx: bid_action.ExpireBidTx{},
		Msg: &bid_action.ExpireBid{},
	}
	counterOffer := common.ExtTx{
		Tx: bid_action.CounterOfferTx{},
		Msg: &bid_action.CounterOffer{},
	}
	bidderDecision := common.ExtTx{
		Tx: bid_action.BidderDecisionTx{},
		Msg: &bid_action.BidderDecision{},
	}
	ownerDecision := common.ExtTx{
		Tx: bid_action.OwnerDecisionTx{},
		Msg: &bid_action.OwnerDecision{},
	}
	appData.ExtTxs = append(appData.ExtTxs, bidCreate)
	appData.ExtTxs = append(appData.ExtTxs, bidCancel)
	appData.ExtTxs = append(appData.ExtTxs, bidExpire)
	appData.ExtTxs = append(appData.ExtTxs, counterOffer)
	appData.ExtTxs = append(appData.ExtTxs, bidderDecision)
	appData.ExtTxs = append(appData.ExtTxs, ownerDecision)

	//load stores
	if dupName, ok := appData.ExtStores["bidMaster"]; ok {
		logger.Errorf("Trying to register external store %s failed, same name already exists", dupName)
		return
	} else {
		appData.ExtStores["bidMaster"] = bid_data.NewBidMasterStore(appData.ChainState)
	}

	//load services
	//first to grab stores from chainstate
	balances := balance.NewStore("b", storage.NewState(appData.ChainState))
	//todo follow this and add any store that is needed from chainstate, this is copying existing stuff to service
	//todo and use ctx.ExtServiceMap in map.go to combine two maps
	appData.ExtServiceMap["bid_query"] = bid_rpc_query.NewService()

	//load beginner and ender functions
	//todo just like
}
