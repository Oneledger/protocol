package app

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"runtime/debug"
	"sort"
	"strconv"

	"github.com/tendermint/tendermint/libs/kv"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"

	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/network_delegation"
	"github.com/Oneledger/protocol/external_apps/common"

	"github.com/Oneledger/protocol/data/balance"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	abciTypes "github.com/tendermint/tendermint/abci/types"

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

		forkMap, err := app.genesisDoc.ForkParams.ToMap()
		if err != nil {
			app.logger.Error("Failed to read fork map", "err", err)
			return ResponseInitChain{}
		}
		utils.PrintStringMap(forkMap, "%s starts at height: %v\n", true)

		err = app.setupState(req.AppStateBytes)
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

func (app *App) applyUpdate(req RequestBeginBlock) error {
	height := req.Header.GetHeight()

	if app.genesisDoc.ForkParams.IsFrankensteinBlock(height) {
		govern := app.Context.govern.WithState(app.Context.deliver)

		options, err := govern.GetStakingOptions()
		if err != nil {
			return err
		}

		options.TopValidatorCount = int64(64)
		options.MinSelfDelegationAmount = *balance.NewAmountFromInt(500_000)

		app.logger.Info("Updating top validator count to", options.TopValidatorCount)
		app.logger.Info("Updating min self delegation amount to", options.MinSelfDelegationAmount)

		err = govern.WithHeight(height).SetStakingOptions(*options)
		if err != nil {
			return errors.Wrap(err, "Setup Staking Options")
		}
		err = govern.WithHeight(height).SetLUH(governance.LAST_UPDATE_HEIGHT_STAKING)
		if err != nil {
			return errors.Wrap(err, "Unable to set last Update height")
		}

		app.logger.Info("Frankenstein applied at block", height)
	}

	// Update last block height and hash
	if app.genesisDoc.ForkParams.IsFrankensteinUpdate(req.Header.GetHeight()) {
		app.Context.stateDB.SetBlockHash(ethcmn.BytesToHash(req.GetHash()))
	}

	return nil
}

func (app *App) blockBeginner() blockBeginner {
	return func(req RequestBeginBlock) ResponseBeginBlock {
		defer app.handlePanic()

		gc := app.getGasCalculator()
		app.Context.deliver = storage.NewState(app.Context.chainstate).WithGas(gc)

		// Apply update at specific height
		if err := app.applyUpdate(req); err != nil {
			panic(err)
		}

		feeOpt, err := app.Context.govern.GetFeeOption()
		if err != nil {
			app.logger.Error("failed to get feeOption", err)
		}
		app.Context.feePool.SetupOpt(feeOpt)

		err = ManageVotes(&req, &app.Context, app.logger)
		if err != nil {
			app.logger.Error("manage votes error", err)
		}

		// update the validator set
		err = app.Context.validators.WithState(app.Context.deliver).Setup(req, app.Context.node.ValidatorAddress())
		if err != nil {
			app.logger.Error("validator set with error", err)
		}

		result := ResponseBeginBlock{
			Events: []abciTypes.Event{},
		}
		// Mature Pending undelegates to delegator's balance
		delegEvent, anyMatured := addMaturedAmountsToBalance(&app.Context, app.logger, &req)
		if anyMatured {
			result.Events = append(result.Events, delegEvent)
		}

		// update malicious list
		err = app.Context.validators.CheckMaliciousValidators(
			app.Context.evidenceStore.WithState(app.Context.deliver),
			app.Context.govern.WithState(app.Context.deliver),
		)
		if err != nil {
			app.logger.Error("malicious set with error", err)
		}

		// update Block Rewards
		blockRewardEvent := handleBlockRewards(&app.Context, req, app.logger)
		result.Events = append(result.Events, blockRewardEvent)

		//update the header to current block
		app.header = req.Header
		//Adds proposals that meet the requirements to either Expired or Finalizing Keys from transaction store
		//Transaction store is not part of chainstate ,it just maintains a list of proposals from BlockBeginner to BlockEnder .Gets cleared at each Block Ender
		AddInternalTX(app.Context.proposalMaster, app.Context.node.ValidatorAddress(), app.header.Height, app.Context.transaction, app.logger)
		functionList, err := app.Context.extFunctions.Iterate(common.BlockBeginner)
		functionParam := common.ExtParam{
			InternalTxStore: app.Context.transaction,
			Logger:          app.logger,
			ActionCtx:       *app.Context.Action(&app.header, app.Context.deliver),
			Validator:       app.Context.node.ValidatorAddress(),
			Header:          app.header,
			Deliver:         app.Context.deliver,
		}
		if err == nil {
			for _, function := range functionList {
				function(functionParam)
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
			loginfo := fmt.Sprintf("checkTx duplicated tx: %s", hex.EncodeToString(utils.GetTransactionHash(msg.Tx)))
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

		feeOk, feeResponse := handler.ProcessFee(txCtx, *tx, gas, storage.Gas(len(msg.Tx)), storage.Gas(response.GasUsed))

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

		txHashBytes := utils.GetTransactionHash(msg.Tx)
		app.Context.stateDB.Prepare(ethcmn.BytesToHash(txHashBytes))

		if cachedResponse, found := app.GetTxFromCache(txHashBytes); found {
			app.logger.Detailf("deliverTx duplicated tx: %s\n", ethcmn.Bytes2Hex(txHashBytes))
			app.Context.stateDB.Finality(cachedResponse.Events)
			return cachedResponse
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
		feeOk, feeResponse := handler.ProcessFee(txCtx, *tx, gas, storage.Gas(len(msg.Tx)), storage.Gas(response.GasUsed))

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

		app.Context.stateDB.Finality(response.Events)

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

		events := app.Context.validators.WithState(app.Context.deliver).GetEvents()
		app.logger.Detailf("Sending events with nodes to tendermint: %+v\n", events)

		app.Context.validators.WithState(app.Context.deliver).ClearEvents()

		ethTrackerlog := log.NewLoggerWithPrefix(app.Context.logWriter, "ethtracker").WithLevel(log.Level(app.Context.cfg.Node.LogLevel))
		doTransitions(app.Context.jobStore, app.Context.btcTrackers.WithState(app.Context.deliver), app.Context.validators)
		doEthTransitions(app.Context.jobStore, app.Context.ethTrackers, app.Context.node.ValidatorAddress(), ethTrackerlog, app.Context.witnesses, app.Context.deliver)
		// Proposals currently in store are cleared if deliver is successful
		// If Expire or Finalize TX returns false,they will added to the proposals queue in the next block
		// Errors are logged at the function level
		// These functions iterate the transactions store
		ExpireProposals(&app.header, &app.Context, app.logger)
		FinalizeProposals(&app.header, &app.Context, app.logger)
		functionList, err := app.Context.extFunctions.Iterate(common.BlockEnder)
		functionParam := common.ExtParam{
			InternalTxStore: app.Context.transaction,
			Logger:          app.logger,
			ActionCtx:       *app.Context.Action(&app.header, app.Context.deliver),
			Validator:       app.Context.node.ValidatorAddress(),
			Header:          app.header,
			Deliver:         app.Context.deliver,
		}
		if err == nil {
			for _, function := range functionList {
				function(functionParam)
			}
		}

		if app.genesisDoc.ForkParams.IsFrankensteinUpdate(req.GetHeight()) {
			// getting bloom if exist
			bloomEvt := app.Context.stateDB.GetBloomEvent()
			if bloomEvt != nil {
				events = append(events, *bloomEvt)
			}
			// Reset all cache after account data has been committed, that make sure node state consistent
			app.Context.stateDB.Reset()
		}

		result := ResponseEndBlock{
			ValidatorUpdates: updates,
			Events:           events,
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
		gc := app.getGasCalculator()
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

func (app *App) getGasCalculator() storage.GasCalculator {
	limit := int64(0)
	if app.genesisDoc.ConsensusParams != nil {
		limit = app.genesisDoc.ConsensusParams.Block.MaxGas
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

			logger.Debug("Processing Tracker : ", t.Type.String(), " | Tracker Name ", t.TrackerName.String(), " | State :", t.State.String(), " | Finality Votes :", t.FinalityVotes)
			_, err := event.EthLockEngine.Process(t.NextStep(), ctx, transition.Status(t.State))
			if err != nil {
				logger.Error("failed to process eth tracker ProcessTypeLock", err)
				continue
			}

		} else if t.Type == ethereum.ProcessTypeRedeem || t.Type == ethereum.ProcessTypeRedeemERC {
			logger.Debug("Processing Tracker : ", t.Type.String(), " | Tracker Name ", t.TrackerName.String(), " | State :", t.State.String(), " | Finality Votes :", t.FinalityVotes)
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
func getRewardForValidator(totalPower *big.Int, validatorPower *big.Int, totalRewards *balance.Amount) *balance.Amount {
	numerator := big.NewInt(0).Mul(totalRewards.BigInt(), validatorPower)
	reward := balance.NewAmountFromBigInt(big.NewInt(0).Div(numerator, totalPower))
	return reward
}

func handleDelegationRewards(delegCtx *network_delegation.DelegationRewardCtx, appCtx *context, kvMap map[string]kv.Pair,
) (resp *network_delegation.DelegationRewardResponse) {
	var err error
	networkDelegators := appCtx.netwkDelegators.WithState(appCtx.deliver)
	//rewardMaster := appCtx.rewardMaster.WithState(appCtx.deliver)
	resp = &network_delegation.DelegationRewardResponse{
		DelegationRewards: balance.NewAmount(0),
		ProposerReward:    balance.NewAmount(0),
		Commission:        balance.NewAmount(0),
	}

	//Get Total Rewards "T" for Delegator Pool
	numerator := big.NewInt(0).Mul(delegCtx.TotalRewards.BigInt(), delegCtx.DelegationPower)
	delegationRewards := balance.NewAmountFromBigInt(big.NewInt(0).Div(numerator, delegCtx.TotalPower))

	//Cut X% from Total Delegation Rewards "T" as Commission "C"
	numerator = big.NewInt(0).Mul(big.NewInt(network_delegation.COMMISSION_PERCENTAGE), delegationRewards.BigInt())
	commission := balance.NewAmountFromBigInt(big.NewInt(0).Div(numerator, big.NewInt(100)))

	//Distribute Y% of Commission "C" to block proposer
	numerator = big.NewInt(0).Mul(big.NewInt(network_delegation.BLOCK_PROPOSER_COMMISSION), commission.BigInt())
	resp.ProposerReward = balance.NewAmountFromBigInt(big.NewInt(0).Div(numerator, big.NewInt(100)))

	//Deduct Commission from Total delegation Rewards
	resp.DelegationRewards, err = delegationRewards.Minus(*commission)
	if err != nil {
		return
	}

	//Deduct Proposer Reward from Commission Amount
	resp.Commission, err = commission.Minus(*resp.ProposerReward)
	if err != nil {
		return
	}

	//Distribute Rewards to Each Delegator
	networkDelegators.Deleg.IterateActiveAmounts(func(addr *keys.Address, coin *balance.Coin) bool {
		//Calculate reward portion for each delegator based on delegated amount
		numerator := big.NewInt(0).Mul(resp.DelegationRewards.BigInt(), coin.Amount.BigInt())
		delegatorReward := balance.NewAmountFromBigInt(big.NewInt(0).Div(numerator, delegCtx.DelegationPower))

		//Add reward to address
		err := networkDelegators.Rewards.AddRewardsBalance(*addr, delegatorReward)
		if err != nil {
			return true
		}
		return false
	})

	//Create Event for Proposer Reward
	proposerKey := "proposer_" + delegCtx.ProposerAddress.String()
	kvMap[proposerKey] = kv.Pair{
		Key:   []byte(proposerKey),
		Value: []byte(resp.ProposerReward.String()),
	}

	//Create Event for Delegation Rewards
	poolList, _ := appCtx.govern.GetPoolList()
	kvMap[poolList["DelegationPool"].String()] = kv.Pair{
		Key:   []byte(poolList["DelegationPool"].String()),
		Value: []byte(resp.DelegationRewards.String()),
	}
	return
}

func handleBlockRewards(appCtx *context, block RequestBeginBlock, logger *log.Logger) abciTypes.Event {
	votes := block.LastCommitInfo.Votes
	lastHeight := block.GetHeader().Height
	rewardMaster := appCtx.rewardMaster.WithState(appCtx.deliver)
	options := rewardMaster.Reward.GetOptions()

	curr, ok := appCtx.currencies.GetCurrencyById(0)
	if !ok {
		return abciTypes.Event{}
	}

	heightKey := "height"
	//Initialize Event for Block Response
	result := abciTypes.Event{}
	kvMap := make(map[string]kv.Pair)

	kvMap[heightKey] = kv.Pair{
		Key:   []byte(heightKey),
		Value: []byte(strconv.FormatInt(lastHeight, 10)),
	}

	//Initialize kvMap with validators containing 0 rewards
	appCtx.validators.Iterate(func(addr keys.Address, validator *identity.Validator) bool {
		kvMap[addr.String()] = kv.Pair{
			Key:   []byte(addr.String()),
			Value: []byte(balance.NewAmount(0).String()),
		}
		return false
	})

	//get total power of active validators
	totValPower := big.NewInt(0)
	validatorPowerMap := make(map[string]*big.Int)
	for _, vote := range votes {
		powerStr := utils.PadZero(strconv.FormatInt(vote.Validator.Power, 10))
		validatorPower, err := balance.NewAmountFromString(powerStr, 10)
		if err != nil {
			return abciTypes.Event{}
		}

		totValPower.Add(totValPower, validatorPower.BigInt())
		validatorPowerMap[keys.Address(vote.Validator.Address).String()] = validatorPower.BigInt()
	}

	//Initialize total Power as total Validator Power
	totalPower := totValPower

	//Add Delegator Pool Balance to the Total Power
	poolList, err := appCtx.govern.GetPoolList()
	if err != nil {
		return abciTypes.Event{}
	}
	delegationPoolCoin, err := appCtx.balances.WithState(appCtx.deliver).GetBalanceForCurr(poolList["DelegationPool"], &curr)
	if err != nil {
		return abciTypes.Event{}
	}
	delegationPower := delegationPoolCoin.Amount.BigInt()
	totalPower.Add(totalPower, delegationPower)

	//get total rewards for the block
	rewardPoolCoin, err := appCtx.balances.WithState(appCtx.deliver).GetBalanceForCurr(poolList["RewardsPool"], &curr)
	if err != nil {
		return abciTypes.Event{}
	}
	totalRewards, err := rewardMaster.RewardCm.PullRewards(lastHeight, rewardPoolCoin.Amount)
	if err != nil {
		return abciTypes.Event{}
	}

	totalConsumed := balance.NewAmount(0)
	delegationResp := &network_delegation.DelegationRewardResponse{}
	if delegationPower.Cmp(big.NewInt(0)) > 0 {
		delegationCtx := &network_delegation.DelegationRewardCtx{
			TotalRewards:    totalRewards,
			DelegationPower: delegationPower,
			TotalPower:      totalPower,
			Height:          lastHeight,
			ProposerAddress: keys.Address(block.Header.ProposerAddress),
		}
		delegationResp = handleDelegationRewards(delegationCtx, appCtx, kvMap)

		//Update Consumed Amount
		totalConsumed = totalConsumed.Plus(*delegationResp.DelegationRewards)
	}

	//Loop through all validators that participated in signing the last block
	for _, vote := range votes {
		//Verify Validator Address
		valAddress := keys.Address(vote.Validator.Address)
		if valAddress.Err() != nil {
			continue
		}
		val, err := appCtx.validators.Get(valAddress)
		if err != nil || len(val.Bytes()) == 0 {
			continue
		}
		if vote.GetSignedLastBlock() {
			//Get Commission and Reward Amounts for Validator
			rewardAmount := getRewardForValidator(totalPower, validatorPowerMap[valAddress.String()], totalRewards)
			commissionAmount := balance.NewAmount(0)
			if delegationPower.Cmp(big.NewInt(0)) > 0 {
				commissionAmount = getRewardForValidator(totValPower, validatorPowerMap[valAddress.String()], delegationResp.Commission)

				if valAddress.String() == keys.Address(block.Header.ProposerAddress).String() {
					commissionAmount = commissionAmount.Plus(*delegationResp.ProposerReward)
				}
			}
			//Add Commission from Delegation rewards
			amount := rewardAmount.Plus(*commissionAmount)
			//Add Amount to reward store
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

	//Update Validators' Matured Amount
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

	//Update delegators' Matured Amount
	matureDelegationRewards(appCtx, &block, kvMap, logger)

	//pass total consumed amount to cumulative db
	_ = rewardMaster.RewardCm.ConsumeRewards(totalConsumed)

	//Populate Event with validator rewards
	result.Type = "block_rewards"

	kvKeys := make([]string, 0, len(kvMap))
	for k := range kvMap {
		kvKeys = append(kvKeys, k)
	}
	sort.Strings(kvKeys)
	for _, key := range kvKeys {
		result.Attributes = append(result.Attributes, kvMap[key])
	}

	return result
}

func (app *App) GetTxFromCache(hash []byte) (abciTypes.ResponseDeliverTx, bool) {
	tx, err := tmrpccore.Tx(nil, hash, false)
	app.logger.Debugf("Got reply for exist by tx hash: %s, err: %s\n", ethcmn.Bytes2Hex(hash), err)
	if tx != nil && tx.Height > 0 {
		return tx.TxResult, true
	}
	return abciTypes.ResponseDeliverTx{}, false
}

func (app *App) VerifyCache(tx []byte) bool {
	reply, err := tmrpccore.Tx(nil, utils.GetTransactionHash(tx), false)
	app.logger.Debugf("Got reply for exist tx: %+v, err: %s\n", reply, err)
	if reply != nil && reply.Height > 0 {
		return true
	}
	return false
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

func ManageVotes(req *RequestBeginBlock, ctx *context, logger *log.Logger) error {
	eopts, err := ctx.govern.WithState(ctx.deliver).GetEvidenceOptions()
	if err != nil {
		logger.Error("error in GetEvidenceOptions")
		return err
	}
	err = ctx.evidenceStore.WithState(ctx.deliver).SetVoteBlock(req.Header.GetHeight(), req.LastCommitInfo.Votes)
	if err != nil {
		logger.Error("error in SetVoteBlock")
		return err
	}
	cv, err := ctx.evidenceStore.WithState(ctx.deliver).GetCumulativeVote()
	if err != nil {
		logger.Error("error in GetCumulativeVote")
		return err
	}
	err = ctx.evidenceStore.WithState(ctx.deliver).SetCumulativeVote(cv, req.Header.GetHeight(), eopts.BlockVotesDiff)
	if err != nil {
		logger.Error("error in SetCumulativeVote")
		return err
	}
	return nil
}

func addMaturedAmountsToBalance(ctx *context, logger *log.Logger, req *RequestBeginBlock) (event abciTypes.Event, any bool) {
	height := req.Header.Height
	delegStore := ctx.netwkDelegators.Deleg.WithState(ctx.deliver)
	balanceStore := ctx.balances.WithState(ctx.deliver)
	c, ok := ctx.currencies.GetCurrencyByName("OLT")
	if !ok {
		logger.Errorf("failed to get OLT as currency from context")
		panic("failed to get OLT as currency from context")
	}
	event = abciTypes.Event{}
	event.Type = "deleg_undelegate"
	event.Attributes = append(event.Attributes, kv.Pair{
		Key:   []byte("height"),
		Value: []byte(strconv.FormatInt(height, 10)),
	})
	// put all the pending amounts at this height directly to delegator's balance
	delegStore.IteratePendingAmounts(height, func(addr *keys.Address, coin *balance.Coin) bool {
		//Add each of them to user's address
		err := balanceStore.AddToAddress(*addr, *coin)
		if err != nil {
			logger.Errorf("failed to add pending undelegation amount at height: %d to address: %s", height, addr.String())
			panic(err)
		}
		//Clear the pending amount
		zeroCoin := c.NewCoinFromAmount(*balance.NewAmount(0))
		err = delegStore.SetPendingAmount(*addr, height, &zeroCoin)
		if err != nil {
			logger.Errorf("failed to clear pending undelegation amount at height: %d for address: %s", height, addr.String())
			panic(err)
		}
		event.Attributes = append(event.Attributes, kv.Pair{
			Key:   []byte(addr.String()),
			Value: []byte(coin.String()),
		})
		return false
	})
	any = len(event.Attributes) > 1
	return
}

func matureDelegationRewards(appCtx *context, req *RequestBeginBlock, kvMap map[string]kv.Pair, logger *log.Logger) {
	networkDelegators := appCtx.netwkDelegators.WithState(appCtx.deliver)
	//Mature pending rewards to delegator's balance
	height := req.Header.Height
	rewardsStore := networkDelegators.Rewards
	balanceStore := appCtx.balances.WithState(appCtx.deliver)
	c, ok := appCtx.currencies.GetCurrencyByName("OLT")
	if !ok {
		logger.Errorf("failed to get OLT as currency from context")
		panic("failed to get OLT as currency from context")
	}
	// put all the pending rewards at the height directly to delegator's balance
	rewardsStore.IteratePD(height, func(delegator keys.Address, amt *balance.Amount) bool {
		//Add each of them to user's address
		coin := c.NewCoinFromAmount(*amt)
		err := balanceStore.AddToAddress(delegator, coin)
		if err != nil {
			logger.Errorf("failed to add pending rewards amount at height: %d to address: %s", height, delegator.String())
			panic(err)
		}
		// clear pending amount
		zero := balance.NewAmount(0)
		err = rewardsStore.SetPendingRewards(delegator, zero, height)
		if err != nil {
			logger.Errorf("failed to clear pending rewards amount at height: %d for address: %s", height, delegator.String())
			panic(err)
		}
		//Create Event for maturing rewards
		rewardsKey := "deleg_rewards_mature_" + delegator.String()
		kvMap[rewardsKey] = kv.Pair{
			Key:   []byte(rewardsKey),
			Value: []byte(coin.String()),
		}
		return false
	})
}
