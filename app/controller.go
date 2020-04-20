package app

import (
	"encoding/hex"
	"fmt"

	"math"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/tendermint/tendermint/types"

	"github.com/Oneledger/protocol/action"
	ceth "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/event"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils/transition"
	"github.com/Oneledger/protocol/version"
)

// The following set of functions will be passed to the abciController

// query connection: for querying the application state; only uses query and Info
func (app *App) infoServer() infoServer {
	return func(info RequestInfo) ResponseInfo {
		defer app.handlePanic()
		//get apphash and height from db
		ver, hash := app.getAppHash()

		result := ResponseInfo{
			Data:             app.name,
			Version:          version.Fullnode.String(),
			LastBlockHeight:  ver,
			LastBlockAppHash: hash,
		}
		app.logger.Info("Server Info:", result, "hash", hex.EncodeToString(hash), "height", ver)
		return result
	}
}

func (app *App) queryer() queryer {
	return func(RequestQuery) ResponseQuery {
		defer app.handlePanic()
		// Do stuff
		return ResponseQuery{}
	}
}

func (app *App) optionSetter() optionSetter {
	return func(req RequestSetOption) ResponseSetOption {
		defer app.handlePanic()
		err := config.Setup(&app.Context.cfg, req.Key, req.Value)
		if err != nil {
			return ResponseSetOption{
				Code: CodeNotOK.uint32(),
				Log:  errors.Wrap(err, "set option").Error(),
			}
		}
		return ResponseSetOption{
			Code: CodeOK.uint32(),
			Info: fmt.Sprintf("set option: key=%s, value=%s ", req.Key, req.Value),
		}
	}
}

// consensus methods: for executing transactions that have been committed. Message sequence is -for every block

func (app *App) chainInitializer() chainInitializer {
	return func(req RequestInitChain) ResponseInitChain {
		defer app.handlePanic()
		app.Context.deliver = storage.NewState(app.Context.chainstate)
		app.Context.govern.WithState(app.Context.deliver)
		app.Context.btcTrackers.WithState(app.Context.deliver)

		err := app.setupState(req.AppStateBytes)
		// This should cause consensus to halt
		if err != nil {
			app.logger.Error("Failed to setupState", "err", err)
			return ResponseInitChain{}
		}

		//update the initial validator set to db, this should always comes after setupState as the currency for
		// validator will be registered by setupState
		validators, err := app.setupValidators(req, app.Context.currencies)
		if err != nil {
			app.logger.Error("Failed to setupValidator", "err", err)
			return ResponseInitChain{}
		}

		app.Context.govern.Initiated()
		app.Context.deliver.Write()

		app.logger.Info("finish chain initialize")
		return ResponseInitChain{Validators: validators}
	}
}

func (app *App) blockBeginner() blockBeginner {
	return func(req RequestBeginBlock) ResponseBeginBlock {
		defer app.handlePanic()
		gc := getGasCalculator(app.genesisDoc.ConsensusParams)
		app.Context.deliver = storage.NewState(app.Context.chainstate).WithGas(gc)

		// update the validator set
		err := app.Context.validators.Setup(req, app.Context.node.ValidatorAddress())
		if err != nil {
			app.logger.Error("validator set with error", err)
		}

		result := ResponseBeginBlock{}

		//update the header to current block
		app.header = req.Header

		app.logger.Detail("Begin Block:", result, "height:", req.Header.Height, "AppHash:", hex.EncodeToString(req.Header.AppHash))
		return result
	}
}

// mempool connection: for checking if transactions should be relayed before they are committed
func (app *App) txChecker() txChecker {
	return func(msg RequestCheckTx) ResponseCheckTx {
		defer app.handlePanic()

		app.Context.check.BeginTxSession()

		tx := &action.SignedTx{}

		err := serialize.GetSerializer(serialize.NETWORK).Deserialize(msg.Tx, tx)
		if err != nil {
			app.logger.Errorf("checkTx failed to deserialize msg: %s, error: %s ", msg, err)
		}

		txCtx := app.Context.Action(&app.header, app.Context.check)
		handler := txCtx.Router.Handler(tx.Type)

		gas := txCtx.State.ConsumedGas()

		ok, err := handler.Validate(txCtx, *tx)
		if err != nil {
			app.logger.Debug("Check Tx invalid: ", err.Error())
			return ResponseCheckTx{
				Code: getCode(ok).uint32(),
				Log:  err.Error(),
			}
		}
		ok, response := handler.ProcessCheck(txCtx, tx.RawTx)

		feeOk, feeResponse := handler.ProcessFee(txCtx, *tx, gas, storage.Gas(len(msg.Tx)))

		result := ResponseCheckTx{
			Code:      getCode(ok && feeOk).uint32(),
			Data:      response.Data,
			Log:       response.Log + feeResponse.Log,
			Info:      response.Info,
			GasWanted: feeResponse.GasWanted,
			GasUsed:   feeResponse.GasUsed,
			Events:    response.Events,
			Codespace: "",
		}

		if !(ok && feeOk) {
			app.Context.check.DiscardTxSession()
		} else {
			app.Context.check.CommitTxSession()
		}

		app.logger.Detail("Check Tx: ", result, "log", response.Log)
		return result

	}
}

func (app *App) txDeliverer() txDeliverer {
	return func(msg RequestDeliverTx) ResponseDeliverTx {
		defer app.handlePanic()

		app.Context.deliver.BeginTxSession()

		tx := &action.SignedTx{}

		err := serialize.GetSerializer(serialize.NETWORK).Deserialize(msg.Tx, tx)
		if err != nil {
			app.logger.Errorf("deliverTx failed to deserialize msg: %s, error: %s ", msg, err)
		}
		txCtx := app.Context.Action(&app.header, app.Context.deliver)

		handler := txCtx.Router.Handler(tx.Type)

		gas := txCtx.State.ConsumedGas()

		ok, response := handler.ProcessDeliver(txCtx, tx.RawTx)

		feeOk, feeResponse := handler.ProcessFee(txCtx, *tx, gas, storage.Gas(len(msg.Tx)))

		result := ResponseDeliverTx{
			Code:      getCode(ok && feeOk).uint32(),
			Data:      response.Data,
			Log:       response.Log + feeResponse.Log,
			Info:      response.Info,
			GasWanted: feeResponse.GasWanted,
			GasUsed:   feeResponse.GasUsed,
			Events:    response.Events,
			Codespace: "",
		}
		app.logger.Detail("Deliver Tx: ", result)

		if !(ok && feeOk) {
			app.Context.deliver.DiscardTxSession()
		} else {
			app.Context.deliver.CommitTxSession()
		}

		return result
	}
}

func (app *App) blockEnder() blockEnder {
	return func(req RequestEndBlock) ResponseEndBlock {
		defer app.handlePanic()

		fee, err := app.Context.feePool.WithState(app.Context.deliver).Get([]byte(fees.POOL_KEY))
		app.logger.Detail("endblock fee", fee, err)
		updates := app.Context.validators.GetEndBlockUpdate(app.Context.ValidatorCtx(), req)
		result := ResponseEndBlock{
			ValidatorUpdates: updates,
			//Tags:             []kv.Pair(nil),
		}
		ethTrackerlog := log.NewLoggerWithPrefix(app.Context.logWriter, "ethtracker").WithLevel(log.Level(app.Context.cfg.Node.LogLevel))
		doTransitions(app.Context.jobStore, app.Context.btcTrackers.WithState(app.Context.deliver), app.Context.validators)
		doEthTransitions(app.Context.jobStore, app.Context.ethTrackers, app.Context.node.ValidatorAddress(), ethTrackerlog, app.Context.witnesses, app.Context.deliver)

		app.logger.Detail("End Block: ", result, "height:", req.Height)

		return result
	}
}

func (app *App) commitor() commitor {
	return func() ResponseCommit {
		defer app.handlePanic()

		hash, ver := app.Context.deliver.Commit()
		app.logger.Detailf("Committed New Block height[%d], hash[%s], versions[%d]", app.header.Height, hex.EncodeToString(hash), ver)

		// update check state by deliver state
		gc := getGasCalculator(app.genesisDoc.ConsensusParams)
		app.Context.check = storage.NewState(app.Context.chainstate).WithGas(gc)
		result := ResponseCommit{
			Data: hash,
		}

		app.logger.Detail("Commit Result", result)
		return result
	}
}

func getCode(ok bool) (code Code) {
	if ok {
		code = CodeOK
	} else {
		code = CodeNotOK
	}
	return
}

func (app *App) getAppHash() (version int64, hash []byte) {
	hash, version = app.Context.chainstate.Hash, app.Context.chainstate.Version
	return
}

func (app *App) handlePanic() {
	if r := recover(); r != nil {
		fmt.Println("panic in controller: ", r)
		debug.PrintStack()
		app.Close()
	}
}

func getGasCalculator(params *types.ConsensusParams) storage.GasCalculator {
	limit := int64(0)
	if params != nil {
		limit = params.Block.MaxGas
	}
	gas := storage.Gas(0)
	if limit < 0 {
		gas = math.MaxInt64
	} else {
		gas = storage.Gas(limit)
	}
	return storage.NewGasCalculator(gas)
}

func doTransitions(js *jobs.JobStore, ts *bitcoin.TrackerStore, validators *identity.ValidatorStore) {

	btcTracker := []bitcoin.Tracker{}
	if js != nil {
		ts.Iterate(func(k, v []byte) bool {

			szlr := serialize.GetSerializer(serialize.PERSISTENT)

			d := &bitcoin.Tracker{}
			err := szlr.Deserialize(v, d)
			if err != nil {
				return false
			}

			btcTracker = append(btcTracker, *d)
			return false
		})
	}

	for _, t := range btcTracker {

		ctx := bitcoin.BTCTransitionContext{&t, js.WithChain(chain.BITCOIN), validators}

		stt, err := event.BtcEngine.Process(t.NextStep(), ctx, transition.Status(t.State))
		if err != nil {
			continue
		}
		if stt != -1 {
			t.State = bitcoin.TrackerState(stt)
			err = ts.SetTracker(t.Name, &t)
		}
	}
}

func doEthTransitions(js *jobs.JobStore, ts *ethereum.TrackerStore, myValAddr keys.Address, logger *log.Logger, witnesses *identity.WitnessStore, deliver *storage.State) {
	ts = ts.WithState(deliver)
	tnames := make([]*ceth.TrackerName, 0, 20)
	ts.WithPrefixType(ethereum.PrefixOngoing).Iterate(func(name *ceth.TrackerName, tracker *ethereum.Tracker) bool {
		tnames = append(tnames, name)
		return false
	})
	for _, name := range tnames {
		deliver.DiscardTxSession()
		deliver.BeginTxSession()
		t, _ := ts.WithPrefixType(ethereum.PrefixOngoing).Get(*name)
		state := t.State
		ctx := ethereum.NewTrackerCtx(t, myValAddr, js.WithChain(chain.ETHEREUM), ts, witnesses, logger)

		if t.Type == ethereum.ProcessTypeLock || t.Type == ethereum.ProcessTypeLockERC {

			logger.Debug("Processing Tracker : ", t.Type.String(), " | State :", t.State.String())
			_, err := event.EthLockEngine.Process(t.NextStep(), ctx, transition.Status(t.State))
			if err != nil {
				logger.Error("failed to process eth tracker ProcessTypeLock", err)
				continue
			}

		} else if t.Type == ethereum.ProcessTypeRedeem || t.Type == ethereum.ProcessTypeRedeemERC {
			logger.Debug("Processing Tracker : ", t.Type.String(), " | State :", t.State.String())
			_, err := event.EthRedeemEngine.Process(t.NextStep(), ctx, transition.Status(t.State))
			if err != nil {
				logger.Error("failed to process eth tracker ProcessTypeRedeem", err)
				continue
			}
		}
		// only set back to chainstate when transition happened.
		if ctx.Tracker.State < 5 && state != ctx.Tracker.State {
			err := ts.WithPrefixType(ethereum.PrefixOngoing).Set(ctx.Tracker)
			if err != nil {
				logger.Error("failed to save eth tracker", err, ctx.Tracker)
				panic(err)
			}
		}
		deliver.CommitTxSession()
	}

}
