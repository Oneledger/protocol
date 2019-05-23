package app

import (
	"encoding/hex"
	"github.com/Oneledger/protocol/log"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/version"
	"github.com/tendermint/tendermint/libs/common"
)

// The following set of functions will be passed to the abciController

// query connection: for querying the application state; only uses Query and Info
func (app *App) infoServer() infoServer {
	return func(info RequestInfo) ResponseInfo {
		return ResponseInfo{
			Data:             app.name,
			Version:          version.Fullnode.String(),
			LastBlockHeight:  app.header.Height,
			LastBlockAppHash: app.header.AppHash,
		}
	}
}

func (app *App) queryer() queryer {
	return func(RequestQuery) ResponseQuery {
		// Do stuff
		return ResponseQuery{}
	}
}

func (app *App) optionSetter() optionSetter {
	return func(RequestSetOption) ResponseSetOption {
		// TODO
		return ResponseSetOption{
			Code: CodeOK.uint32(),
		}
	}
}

// consensus methods: for executing transactions that have been committed. Message sequence is -for every block

func (app *App) chainInitializer() chainInitializer {
	return func(req RequestInitChain) ResponseInitChain {
		err := app.setupState(req.AppStateBytes)
		// This should cause consensus to halt
		if err != nil {
			app.logger.Error("Failed to setupState", "err", err)
			return ResponseInitChain{}
		}
		return ResponseInitChain{}
	}
}

var startTime time.Time
var endTime time.Time
var startTx int64
var endTx int64

func (app *App) blockBeginner() blockBeginner {
	return func(req RequestBeginBlock) ResponseBeginBlock {

		//update the validator set
		//err := app.Context.validators.Set(req)
		//if err != nil {
		//	app.logger.Error("validator set with error", err)
		//}
		//update the header to current block
		//todo: store the header in persistent db
		app.header = req.Header

		result := ResponseBeginBlock{
			Tags: []common.KVPair(nil),
		}

		if req.Header.Height == 3 {
			startTime = req.Header.Time
			startTx = req.Header.TotalTxs
		}

		if req.Header.Height == 10 {
			endTime = req.Header.Time
			endTx = req.Header.TotalTxs
			loadtest(req.Header, app.logger)
		}

		if req.Header.Height == 30 {
			endTime = req.Header.Time
			endTx = req.Header.TotalTxs
			loadtest(req.Header, app.logger)
		}

		if req.Header.Height == 50 {
			endTime = req.Header.Time
			endTx = req.Header.TotalTxs
			loadtest(req.Header, app.logger)
		}

		app.logger.Debug("Begin Block:", result)
		return result
	}
}

func loadtest(head Header, logger *log.Logger) {
	tps := float64(endTx-startTx) / (endTime.Sub(startTime).Seconds())
	blktime := (endTime.Sub(startTime).Seconds()) / float64(head.Height-3)
	logger.Infof("Loadtest metric height=%d, tx/b=%d, blktime=%3f , tps=%3f", head.Height, head.TotalTxs/head.Height, blktime, tps)
}

// mempool connection: for checking if transactions should be relayed before they are committed
func (app *App) txChecker() txChecker {
	return func(msg []byte) ResponseCheckTx {
		tx := &action.BaseTx{}

		err := serialize.GetSerializer(serialize.NETWORK).Deserialize(msg, tx)
		if err != nil {
			app.logger.Errorf("failed to deserialize msg: %s, error: %s ", msg, err)
		}
		txCtx := app.Context.Action()

		handler := txCtx.Router.Handler(tx.Data)

		ok, response := handler.ProcessCheck(txCtx, tx.Data, tx.Fee)

		var code Code
		if ok {
			code = CodeOK
		} else {
			code = CodeNotOK
		}
		result := ResponseCheckTx{
			Code:      code.uint32(),
			Data:      response.Data,
			Log:       response.Log,
			Info:      response.Info,
			GasWanted: response.GasWanted,
			GasUsed:   response.GasUsed,
			Tags:      response.Tags,
			Codespace: "",
		}
		app.logger.Debug("Check Tx: ", result)
		return result

	}
}

func (app *App) txDeliverer() txDeliverer {
	return func(msg []byte) ResponseDeliverTx {
		tx := &action.BaseTx{}

		err := serialize.GetSerializer(serialize.NETWORK).Deserialize(msg, tx)
		if err != nil {
			app.logger.Errorf("failed to deserialize msg: %s, error: %s ", msg, err)
		}
		txCtx := app.Context.Action()

		handler := txCtx.Router.Handler(tx.Data)

		ok, response := handler.ProcessDeliver(txCtx, tx.Data, tx.Fee)

		var code Code
		if ok {
			code = CodeOK
		} else {
			code = CodeNotOK
		}

		result := ResponseDeliverTx{
			Code:      code.uint32(),
			Data:      response.Data,
			Log:       response.Log,
			Info:      response.Info,
			GasWanted: response.GasWanted,
			GasUsed:   response.GasUsed,
			Tags:      response.Tags,
			Codespace: "",
		}
		app.logger.Debug("Deliver Tx: ", result)
		return result
	}
}

func (app *App) blockEnder() blockEnder {
	return func(req RequestEndBlock) ResponseEndBlock {

		updates := app.Context.validators.GetEndBlockUpdate(app.Context.ValidatorCtx(), req)

		result := ResponseEndBlock{
			ValidatorUpdates: updates,
			Tags:             []common.KVPair(nil),
		}
		app.logger.Debug("End Block: ", result)
		return result
	}
}

func (app *App) commitor() commitor {
	return func() ResponseCommit {

		// Commit any pending changes.
		hash, ver := app.Context.balances.Commit()

		app.logger.Debugf("Committed New Block hash[%s], version[%d]", hex.EncodeToString(hash), ver)

		result := ResponseCommit{
			Data: hash,
		}

		app.logger.Debug("Commit Result", result)
		return result
	}
}
