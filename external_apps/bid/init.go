package bid

import (
	"fmt"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/external_apps/bid/bid_action"
	"github.com/Oneledger/protocol/external_apps/bid/bid_block_func"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
	"github.com/Oneledger/protocol/external_apps/bid/bid_rpc/bid_rpc_query"
	"github.com/Oneledger/protocol/external_apps/bid/bid_rpc/bid_rpc_tx"
	"github.com/Oneledger/protocol/external_apps/common"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"os"
)

var logger *log.Logger

func init() {
	fmt.Println("init from bid/init")
	//common.Handlers.Register(LoadAppData)
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
	balances := balance.NewStore("b", storage.NewState(appData.ChainState))
	domains := ons.NewDomainStore("ons", storage.NewState(appData.ChainState))
	olt := balance.Currency{Id: 0, Name: "OLT", Chain: chain.ONELEDGER, Decimal: 18, Unit: "nue"}
	currencies := balance.NewCurrencySet()
	err := currencies.Register(olt)
	if err != nil {
		logger.Errorf("failed to register currency %s", olt.Name, err)
		return
	}
	logWriter := os.Stdout
	logger := log.NewLoggerWithPrefix(logWriter, "rpc").WithLevel(log.Level(4))
	appData.ExtServiceMap[bid_rpc_query.Name()] = bid_rpc_query.NewService(balances, currencies, domains, logger, bid_data.NewBidMasterStore(appData.ChainState))
	appData.ExtServiceMap[bid_rpc_tx.Name()] = bid_rpc_tx.NewService(balances, logger)

	//load beginner and ender functions
	err = appData.ExtBlockFuncs.Add(common.BlockBeginner, bid_block_func.AddExpireBidTxToQueue)
	if err != nil {
		logger.Errorf("failed to load block beginner func", err)
		return
	}
	err = appData.ExtBlockFuncs.Add(common.BlockEnder, bid_block_func.PopExpireBidTxFromQueue)
	if err != nil {
		logger.Errorf("failed to load block ender func", err)
		return
	}

}
