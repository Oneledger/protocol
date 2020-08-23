package app

import (
	"github.com/Oneledger/protocol/action"
	bid_action "github.com/Oneledger/protocol/action/bidding"
	"github.com/Oneledger/protocol/data/bidding"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/transactions"
	"github.com/google/uuid"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"time"
)

// Follow internalTx.go as a sample
// Function for block Beginner
func AddExpireBidTxToQueue(i interface{}) {
	// Add a store similar to the transaction store to external stores .
	// Access that store through app.Context.extStores.
	// Add transaction to the queue from there .

	// 1. get all the needed stores
	app, ok := i.(*App)
	if ok == false {
		app.logger.Error("failed to assert app in block beginner")
	}
	bidMaster, err := app.Context.extStores.Get("bidMaster")
	if err != nil {
		app.logger.Error("failed to get bid master store in block beginner", err)
	}
	bidMasterStore, ok := bidMaster.(*bidding.BidMasterStore)
	if ok == false {
		app.logger.Error("failed to assert bid master store in block beginner", err)
	}
	bidConvStore := bidMasterStore.BidConv

	internalBidTx, err := app.Context.extStores.Get("internalBidTx")
	if err != nil {
		app.logger.Error("failed to get internal bid tx store in block beginner", err)
	}
	internalBidTxStore, ok := internalBidTx.(*transactions.TransactionStore)
	if ok == false {
		app.logger.Error("failed to assert internal bid tx store in block beginner", err)
	}

	// 2. iterate all the bid conversations and pick the ones that need to be expired
	bidConvStore.Iterate(func(id bidding.BidConvId, bidConv *bidding.BidConv) bool {
		// check expiry
		deadLine := time.Unix(bidConv.DeadlineUTC, 0)

		if deadLine.Before(app.header.Time) {
			// get tx
			tx, err := GetExpireBidTX(bidConv.BidConvId, app.Context.node.ValidatorAddress())
			if err != nil {
				app.logger.Error("Error in building TX of type RequestDeliverTx(expire)", err)
				return true
			}
			// Add tx to expired prefix of transaction store
			err = internalBidTxStore.AddExpired(string(bidConv.BidConvId), &tx)
			if err != nil {
				app.logger.Error("Error in adding to Expired Queue :", err)
				return true
			}

			// Commit the state
			internalBidTxStore.State.Commit()
		}
		return false
	})
}

func GetExpireBidTX(bidConvId bidding.BidConvId, validatorAddress keys.Address) (abciTypes.RequestDeliverTx, error) {
	expireBid := &bid_action.ExpireBid{
		BidConvId:       bidConvId,
		ValidatorAddress: validatorAddress,
	}

	txData, err := expireBid.Marshal()
	if err != nil {
		return RequestDeliverTx{}, err
	}

	internalFinalizeTx := abciTypes.RequestDeliverTx{
		Tx:                   txData,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	return internalFinalizeTx, nil
}


//Function for block Ender
func PopExpireBidTxFromQueue(i interface{}) {
	//Same as above
	//Pop The TX ,call deliverTX on it
	//Use deliverTxSession to commit or ignore the error

	//1. get the internal bid tx store
	app, ok := i.(*App)
	if ok == false {
		app.logger.Error("failed to assert app in block beginner")
	}
	internalBidTx, err := app.Context.extStores.Get("internalBidTx")
	if err != nil {
		app.logger.Error("failed to get internal bid tx store in block beginner", err)
	}
	internalBidTxStore, ok := internalBidTx.(*transactions.TransactionStore)
	if ok == false {
		app.logger.Error("failed to assert internal bid tx store in block beginner", err)
	}

	//2. get all the pending txs
	var expiredBidConvs []abciTypes.RequestDeliverTx
	internalBidTxStore.IterateExpired(func(key string, tx *abciTypes.RequestDeliverTx) bool {
		expiredBidConvs = append(expiredBidConvs, *tx)
		return false
	})

	//3. execute all the txs
	for _, bidConv := range expiredBidConvs {
		app.Context.deliver.BeginTxSession()
		actionctx := app.Context.Action(&app.header, app.Context.deliver)
		txData := bidConv.Tx
		newExpire := bid_action.ExpireBid{}
		err := newExpire.Unmarshal(txData)
		if err != nil {
			app.logger.Error("Unable to UnMarshal TX(Expire) :", txData)
			continue
		}
		uuidNew, _ := uuid.NewUUID()
		rawTx := action.RawTx{
			Type: action.BID_EXPIRE,
			Data: txData,
			Fee:  action.Fee{},
			Memo: uuidNew.String(),
		}
		ok, _ := newExpire.ProcessDeliver(actionctx, rawTx)
		if !ok {
			app.logger.Error("Failed to Expire : ", txData, "Error : ", err)
			app.Context.deliver.DiscardTxSession()
			continue
		}
		app.Context.deliver.CommitTxSession()
	}

	//4. clear txs in transaction store
	internalBidTxStore.IterateExpired(func(key string, tx *abciTypes.RequestDeliverTx) bool {
		ok, err := internalBidTxStore.DeleteExpired(key)
		if !ok {
			app.logger.Error("Failed to clear expired proposals queue :", err)
			return true
		}
		return false
	})
	internalBidTxStore.State.Commit()
}