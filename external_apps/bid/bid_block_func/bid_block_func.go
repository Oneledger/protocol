package bid_block_func

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/external_apps/bid/bid_action"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
	"github.com/google/uuid"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

// Function for block Beginner
func AddExpireBidTxToQueue(i interface{}) {
	// Add a store similar to the transaction store to external stores .
	// Access that store through app.Context.extStores.
	// Add transaction to the queue from there .

	// 1. get all the needed stores
	bidParam, ok := i.(BidParam)
	if ok == false {
		bidParam.Logger.Error("failed to assert bidParam in block beginner")
		return
	}
	bidMaster, err := bidParam.ActionCtx.ExtStores.Get("bidMaster")
	if err != nil {
		bidParam.Logger.Error("failed to get bid master store in block beginner", err)
		return
	}
	bidMasterStore, ok := bidMaster.(*bid_data.BidMasterStore)
	if ok == false {
		bidParam.Logger.Error("failed to assert bid master store in block beginner", err)
		return
	}

	bidConvStore := bidMasterStore.BidConv


	// 2. iterate all the bid conversations and pick the ones that need to be expired
	bidConvStore.Iterate(func(id bid_data.BidConvId, bidConv *bid_data.BidConv) bool {
		// check expiry
		deadLine := time.Unix(bidConv.DeadlineUTC, 0)

		if deadLine.Before(bidParam.Header.Time) {
			// get tx
			tx, err := GetExpireBidTX(bidConv.BidConvId, bidParam.Validator)
			if err != nil {
				bidParam.Logger.Error("Error in building TX of type RequestDeliverTx(expire)", err)
				return true
			}
			// Add tx to expired prefix of transaction store
			err = bidParam.InternalTxStore.AddCustom("bidExpire", string(bidConv.BidConvId), &tx)
			if err != nil {
				bidParam.Logger.Error("Error in adding to Expired Queue :", err)
				return true
			}

			// Commit the state
			bidParam.InternalTxStore.State.Commit()
		}
		return false
	})
}

func GetExpireBidTX(bidConvId bid_data.BidConvId, validatorAddress keys.Address) (abci.RequestDeliverTx, error) {
	expireBid := &bid_action.ExpireBid{
		BidConvId:       bidConvId,
		ValidatorAddress: validatorAddress,
	}

	txData, err := expireBid.Marshal()
	if err != nil {
		return abci.RequestDeliverTx{}, err
	}

	internalFinalizeTx := abci.RequestDeliverTx{
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
	bidParam, ok := i.(BidParam)
	if ok == false {
		bidParam.Logger.Error("failed to assert bidParam in block ender")
		return
	}

	//2. get all the pending txs
	var expiredBidConvs []abci.RequestDeliverTx
	bidParam.InternalTxStore.IterateCustom("bidExpire", func(key string, tx *abci.RequestDeliverTx) bool {
		expiredBidConvs = append(expiredBidConvs, *tx)
		return false
	})

	//3. execute all the txs
	for _, bidConv := range expiredBidConvs {
		bidParam.Deliver.BeginTxSession()
		actionctx := bidParam.ActionCtx
		txData := bidConv.Tx
		newExpireTx := bid_action.ExpireBidTx{}
		newExpire := bid_action.ExpireBid{}
		err := newExpire.Unmarshal(txData)
		if err != nil {
			bidParam.Logger.Error("Unable to UnMarshal TX(Expire) :", txData)
			continue
		}
		uuidNew, _ := uuid.NewUUID()
		rawTx := action.RawTx{
			Type: bid_action.BID_EXPIRE,
			Data: txData,
			Fee:  action.Fee{},
			Memo: uuidNew.String(),
		}
		ok, _ := newExpireTx.ProcessDeliver(&actionctx, rawTx)
		if !ok {
			bidParam.Logger.Error("Failed to Expire : ", txData, "Error : ", err)
			bidParam.Deliver.DiscardTxSession()
			continue
		}
		bidParam.Deliver.CommitTxSession()
	}

	//4. clear txs in transaction store
	bidParam.InternalTxStore.IterateCustom("bidExpire", func(key string, tx *abci.RequestDeliverTx) bool {
		ok, err := bidParam.InternalTxStore.DeleteCustom("bidExpire", key)
		if !ok {
			bidParam.Logger.Error("Failed to clear expired proposals queue :", err)
			return true
		}
		return false
	})
	bidParam.InternalTxStore.State.Commit()
}
