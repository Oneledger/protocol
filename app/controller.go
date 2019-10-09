package app

import (
	"encoding/hex"
	"math"

	"github.com/Oneledger/protocol/data/fees"

	"github.com/tendermint/tendermint/types"

	"github.com/Oneledger/protocol/storage"

	"github.com/Oneledger/protocol/utils"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/version"
	"github.com/tendermint/tendermint/libs/common"
)

// The following set of functions will be passed to the abciController

// query connection: for querying the application state; only uses query and Info
func (app *App) infoServer() infoServer {
	return func(info RequestInfo) ResponseInfo {

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
		app.Context.deliver = storage.NewState(app.Context.chainstate)
		app.Context.govern.WithState(app.Context.deliver)
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
		gc := getGasCalculator(app.genesisDoc.ConsensusParams)
		app.Context.deliver = storage.NewState(app.Context.chainstate).WithGas(gc)

		// update the validator set
		err := app.Context.validators.Setup(req)
		if err != nil {
			app.logger.Error("validator set with error", err)
		}

		result := ResponseBeginBlock{
			Tags: []common.KVPair(nil),
		}

		//update the header to current block
		app.header = req.Header

		app.logger.Debug("Begin Block:", result, "height:", req.Header.Height, "AppHash:", hex.EncodeToString(req.Header.AppHash))
		return result
	}
}

// mempool connection: for checking if transactions should be relayed before they are committed
func (app *App) txChecker() txChecker {
	return func(msg []byte) ResponseCheckTx {
		tx := &action.SignedTx{}

		err := serialize.GetSerializer(serialize.NETWORK).Deserialize(msg, tx)
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

		feeOk, feeResponse := handler.ProcessFee(txCtx, *tx, gas, storage.Gas(len(msg)))

		result := ResponseCheckTx{
			Code:      getCode(ok && feeOk).uint32(),
			Data:      response.Data,
			Log:       response.Log + feeResponse.Log,
			Info:      response.Info,
			GasWanted: feeResponse.GasWanted,
			GasUsed:   feeResponse.GasUsed,
			Tags:      response.Tags,
			Codespace: "",
		}
		app.logger.Debug("Check Tx: ", result, "log", response.Log)
		return result

	}
}

func (app *App) txDeliverer() txDeliverer {
	return func(msg []byte) ResponseDeliverTx {
		tx := &action.SignedTx{}

		err := serialize.GetSerializer(serialize.NETWORK).Deserialize(msg, tx)
		if err != nil {
			app.logger.Errorf("deliverTx failed to deserialize msg: %s, error: %s ", msg, err)
		}
		txCtx := app.Context.Action(&app.header, app.Context.deliver)

		handler := txCtx.Router.Handler(tx.Type)

		gas := txCtx.State.ConsumedGas()

		ok, response := handler.ProcessDeliver(txCtx, tx.RawTx)

		feeOk, feeResponse := handler.ProcessFee(txCtx, *tx, gas, storage.Gas(len(msg)))

		result := ResponseDeliverTx{
			Code:      getCode(ok && feeOk).uint32(),
			Data:      response.Data,
			Log:       response.Log + feeResponse.Log,
			Info:      response.Info,
			GasWanted: feeResponse.GasWanted,
			GasUsed:   feeResponse.GasUsed,
			Tags:      response.Tags,
			Codespace: "",
		}
		app.logger.Debug("Deliver Tx: ", result)
		return result
	}
}

func (app *App) blockEnder() blockEnder {
	return func(req RequestEndBlock) ResponseEndBlock {

		fee, err := app.Context.feePool.WithState(app.Context.deliver).Get([]byte(fees.POOL_KEY))
		app.logger.Debug("endblock fee", fee, err)
		updates := app.Context.validators.GetEndBlockUpdate(app.Context.ValidatorCtx(), req)
		result := ResponseEndBlock{
			ValidatorUpdates: updates,
			Tags:             []common.KVPair(nil),
		}

		app.logger.Debug("End Block: ", result, "height:", req.Height)
		return result
	}
}

func (app *App) commitor() commitor {
	return func() ResponseCommit {

		// Commit any pending changes.
		hash, ver := app.Context.deliver.Commit()
		app.logger.Debugf("Committed New Block height[%d], hash[%s], versions[%d]", app.header.Height, hex.EncodeToString(hash), ver)

		// update check state by deliver state
		gc := getGasCalculator(app.node.GenesisDoc().ConsensusParams)
		app.Context.check = storage.NewState(app.Context.chainstate).WithGas(gc)
		result := ResponseCommit{
			Data: hash,
		}

		app.logger.Debug("Commit Result", result)
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

// TODO: make appHash to use a commiter function to finish the commit and hashing for all the store that passed
type appHash struct {
	Hashes [][]byte `json:"hashes"`
}

func (ah *appHash) hash() []byte {
	result, _ := serialize.GetSerializer(serialize.JSON).Serialize(ah)
	return utils.Hash(result)
}

func (app *App) getAppHash() (version int64, hash []byte) {
	hash, version = app.Context.chainstate.Hash, app.Context.chainstate.Version
	return
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
