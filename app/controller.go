package app

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"runtime/debug"
	"strconv"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/rewards"

	"github.com/pkg/errors"

	abciTypes "github.com/tendermint/tendermint/abci/types"
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
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils"
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

		// update Block Rewards
		blockRewardEvent := handleBlockRewards(app.Context.validators, app.Context.balances,
			app.Context.rewardMaster.WithState(app.Context.deliver), app.Context.currencies, req)

		result := ResponseBeginBlock{
			Events: []abciTypes.Event{blockRewardEvent},
		}

		//update the header to current block
		app.header = req.Header
		//Adds proposals that meet the requirements to either Expired or Finalizing Keys from transaction store
		//Transaction store is not part of chainstate ,it just maintains a list of proposals from BlockBeginner to BlockEnder .Gets cleared at each Block Ender
		AddInternalTX(app.Context.proposalMaster, app.Context.node.ValidatorAddress(), app.header.Height, app.Context.transaction, app.logger)

		functionList, err := app.Context.controllerFunctions.Iterate(BlockBeginner)
		if err == nil {
			for _, controllerFunction := range functionList {
				controllerFunction.Function(controllerFunction.FunctionParam)
			}
		}
		app.logger.Detail("Begin Block:", result, "height:", req.Header.Height, "AppHash:", hex.EncodeToString(req.Header.AppHash))
		return result
	}
}

// mempool connection: for checking if transactions should be relayed before they are committed
func (app *App) txChecker() txChecker {
	return func(msg RequestCheckTx) ResponseCheckTx {
		defer app.handlePanic()

		if app.VerifyCache(msg.Tx) {
			loginfo := fmt.Sprintf("checkTx duplicated tx: %s", hex.EncodeToString(utils.SHA2(msg.Tx)))
			app.logger.Detail(loginfo)
			return ResponseCheckTx{
				Code: CodeNotOK.uint32(),
				Log:  loginfo,
			}
		}

		app.Context.check.BeginTxSession()

		tx := &action.SignedTx{}

		err := serialize.GetSerializer(serialize.NETWORK).Deserialize(msg.Tx, tx)
		if err != nil {
			app.logger.Errorf("checkTx failed to deserialize msg: %v, error: %s ", msg, err)
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

		logString := marshalLog(ok, response, feeResponse)

		result := ResponseCheckTx{
			Code:      getCode(ok && feeOk).uint32(),
			Data:      response.Data,
			Log:       logString,
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

		if app.VerifyCache(msg.Tx) {
			loginfo := fmt.Sprintf("deliverTx duplicated tx: %s", hex.EncodeToString(utils.SHA2(msg.Tx)))
			app.logger.Detail(loginfo)
			return ResponseDeliverTx{
				Code: CodeNotOK.uint32(),
				Log:  loginfo,
			}
		}

		app.Context.deliver.BeginTxSession()

		tx := &action.SignedTx{}

		err := serialize.GetSerializer(serialize.NETWORK).Deserialize(msg.Tx, tx)
		if err != nil {
			app.logger.Errorf("deliverTx failed to deserialize msg: %v, error: %s ", msg, err)
		}
		txCtx := app.Context.Action(&app.header, app.Context.deliver)

		handler := txCtx.Router.Handler(tx.Type)

		gas := txCtx.State.ConsumedGas()

		ok, response := handler.ProcessDeliver(txCtx, tx.RawTx)

		feeOk, feeResponse := handler.ProcessFee(txCtx, *tx, gas, storage.Gas(len(msg.Tx)))

		logString := marshalLog(ok, response, feeResponse)

		result := ResponseDeliverTx{
			Code:      getCode(ok && feeOk).uint32(),
			Data:      response.Data,
			Log:       logString,
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

		updates := app.Context.validators.WithState(app.Context.deliver).GetEndBlockUpdate(app.Context.ValidatorCtx(), req)
		app.logger.Detailf("Sending updates with nodes to tendermint: %+v\n", updates)

		ethTrackerlog := log.NewLoggerWithPrefix(app.Context.logWriter, "ethtracker").WithLevel(log.Level(app.Context.cfg.Node.LogLevel))
		doTransitions(app.Context.jobStore, app.Context.btcTrackers.WithState(app.Context.deliver), app.Context.validators)
		doEthTransitions(app.Context.jobStore, app.Context.ethTrackers, app.Context.node.ValidatorAddress(), ethTrackerlog, app.Context.witnesses, app.Context.deliver)
		// Proposals currently in store are cleared if deliver is successful
		// If Expire or Finalize TX returns false,they will added to the proposals queue in the next block
		// Errors are logged at the function level
		// These functions iterate the transactions store
		ExpireProposals(&app.header, &app.Context, app.logger)
		FinalizeProposals(&app.header, &app.Context, app.logger)

		functionList, err := app.Context.controllerFunctions.Iterate(BlockEnder)
		if err == nil {
			for _, controllerFunction := range functionList {
				controllerFunction.Function(controllerFunction.FunctionParam)
			}
		}
		result := ResponseEndBlock{
			ValidatorUpdates: updates,
		}

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

//Get Individual Validator reward based on power
func getRewardForValidator(totalPower int64, validatorPower int64, totalRewards *balance.Amount) *balance.Amount {
	numerator := big.NewInt(0).Mul(totalRewards.BigInt(), big.NewInt(validatorPower))
	reward := balance.NewAmountFromBigInt(big.NewInt(0).Div(numerator, big.NewInt(totalPower)))
	return reward
}

func handleBlockRewards(validators *identity.ValidatorStore, balances *balance.Store, rewardMaster *rewards.RewardMasterStore,
	currency *balance.CurrencySet, block RequestBeginBlock) abciTypes.Event {

	votes := block.LastCommitInfo.Votes
	lastHeight := block.GetHeader().Height
	options := rewardMaster.Reward.GetOptions()

	heightKey := "height"

	//Initialize Event for Block Response
	result := abciTypes.Event{}
	kvMap := make(map[string]kv.Pair)

	kvMap[heightKey] = kv.Pair{
		Key:   []byte(heightKey),
		Value: []byte(strconv.FormatInt(lastHeight, 10)),
	}

	//Initialize kvMap with validators containing 0 rewards
	validators.Iterate(func(addr keys.Address, validator *identity.Validator) bool {
		kvMap[addr.String()] = kv.Pair{
			Key:   []byte(addr.String()),
			Value: []byte(balance.NewAmount(0).String()),
		}
		return false
	})

	//get total power of active validators
	totalPower := int64(0)
	for _, vote := range votes {
		totalPower += vote.Validator.Power
	}

	//get total rewards for the block
	rewardPoolBalance, err := balances.GetBalance(keys.Address(options.RewardPoolAddress), currency)
	curr, ok := currency.GetCurrencyById(0)
	if err != nil || !ok {
		return abciTypes.Event{}
	}
	currency.GetCurrencyById(0)
	coin := rewardPoolBalance.GetCoin(curr)
	totalRewards, err := rewardMaster.RewardCm.PullRewards(lastHeight, coin.Amount)
	if err != nil {
		return abciTypes.Event{}
	}

	//Loop through all validators that participated in signing the last block
	totalConsumed := balance.NewAmount(0)
	for _, vote := range votes {
		//Verify Validator Address
		valAddress := keys.Address(vote.Validator.Address)
		if valAddress.Err() != nil {
			continue
		}
		val, err := validators.Get(valAddress)
		if err != nil || len(val.Bytes()) == 0 {
			continue
		}

		if vote.GetSignedLastBlock() {
			amount := getRewardForValidator(totalPower, vote.Validator.Power, totalRewards)
			err = rewardMaster.Reward.AddToAddress(valAddress, lastHeight, amount)
			if err != nil {
				continue
			}
			//Add to Consumed amount
			totalConsumed = totalConsumed.Plus(*amount)

			//Record Amount in kvMap
			kvMap[valAddress.String()] = kv.Pair{
				Key:   []byte(valAddress.String()),
				Value: []byte(amount.String()),
			}
		}
	}

	//Update Matured Amount
	rewardMaster.Reward.IterateAddrList(func(addr keys.Address) bool {
		matured := lastHeight % options.RewardInterval
		if matured == 0 {
			//Add rewards at chunk n - 2 to cumulative store
			maturedAmount, err := rewardMaster.Reward.GetMaturedAmount(addr, lastHeight)
			if err != nil {
				return false
			}
			err = rewardMaster.RewardCm.AddMaturedBalance(addr, maturedAmount)
			if err != nil {
				return false
			}
		}
		return false
	})

	//pass total consumed amount to cumulative db
	_ = rewardMaster.RewardCm.ConsumeRewards(totalConsumed)

	//Populate Event with validator rewards
	result.Type = "block_rewards"
	for _, value := range kvMap {
		result.Attributes = append(result.Attributes, value)
	}

	return result
}

func (app *App) VerifyCache(tx []byte) bool {
	hash := utils.SHA2(tx)
	return app.Context.internalService.ExistTx(hash)
}

func marshalLog(ok bool, response action.Response, feeResponse action.Response) string {
	var errorObj codes.ProtocolError
	var err error
	if response.Log == "" && feeResponse.Log == "" {
		return ""
	}
	if !ok {
		errorObj, err = codes.UnMarshalError(response.Log)
		if err != nil {
			// means response.Log is a regular string, from where error marshal has not
			// been done(will do it later)
			errorObj = codes.ProtocolError{
				Code: codes.GeneralErr,
				Msg:  response.Log,
			}
		}

	}
	if feeResponse.Log != "" {
		errorObj.Msg += ", fee response log: " + feeResponse.Log
	}

	return errorObj.Marshal()

}
