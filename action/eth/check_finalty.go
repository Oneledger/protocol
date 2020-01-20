package eth

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/keys"
)

type ReportFinality struct {
	TrackerName      ethereum.TrackerName
	Locker           action.Address
	ValidatorAddress action.Address
	VoteIndex        int64
	Refund           bool
}

var _ action.Msg = &ReportFinality{}

func (m *ReportFinality) Signers() []action.Address {
	return []action.Address{
		m.ValidatorAddress,
	}
}

func (m *ReportFinality) Type() action.Type {
	return action.ETH_REPORT_FINALITY_MINT
}

func (m *ReportFinality) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.ETH_REPORT_FINALITY_MINT.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: m.Locker.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(m.TrackerName.Hex()),
	}
	tag4 := common.KVPair{
		Key:   []byte("tx.validator"),
		Value: m.ValidatorAddress.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

func (m *ReportFinality) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ReportFinality) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

var _ action.Tx = reportFinalityMintTx{}

type reportFinalityMintTx struct {
}

func (reportFinalityMintTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {

	f := &ReportFinality{}
	err := f.Unmarshal(signedTx.Data)
	if err != nil {
		ctx.Logger.Error(err)
		return false, errors.Wrap(err, action.ErrWrongTxType.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), f.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	if f.VoteIndex < 0 {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (reportFinalityMintTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runCheckFinality(ctx, tx)
}

func (reportFinalityMintTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runCheckFinality(ctx, tx)
}

func runCheckFinality(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	f := &ReportFinality{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.ETHTrackers.Get(f.TrackerName)
	if err != nil {
		ctx.Logger.Error(err, "err getting tracker")
	}

	if tracker.Finalized() {
		return true, action.Response{Log: "Tracker already finalized"}
	}
	err = tracker.AddVote(f.ValidatorAddress, f.VoteIndex, true)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to add vote").Error()}
	}

	if tracker.Finalized() {
		ctx.Logger.Info("Finalizing Tracker [ Minting / Burning ]  | Process Type : ", trackerlib.GetProcessTypeString(tracker.Type))
		if tracker.Type == trackerlib.ProcessTypeLock {
			err := mintTokens(ctx, tracker, *f)
			if err != nil {
				return false, action.Response{Log: errors.Wrap(err, "unable to mint tokens").Error()}
			}
		} else if tracker.Type == trackerlib.ProcessTypeRedeem {
			err := burnTokens(ctx, tracker, *f)
			if err != nil {
				return false, action.Response{Log: errors.Wrap(err, "unable to burn tokens").Error()}
			}
		} else if tracker.Type == trackerlib.ProcessTypeLockERC {
			err := mintERC20tokens(ctx, tracker, *f)
			if err != nil {
				return false, action.Response{Log: errors.Wrap(err, "unable to mint tokens").Error()}
			}
		} else if tracker.Type == trackerlib.ProcessTypeRedeemERC {
			err := burnERC20Tokens(ctx,tracker, *f)
			if err != nil {
				return false, action.Response{Log: errors.Wrap(err, "unable to burn tokens").Error()}
			}
		}

		return true, action.Response{Log: "Operation successful"}
	}

	err = ctx.ETHTrackers.Set(tracker)
	if err != nil {
		ctx.Logger.Info("Unable to save the tracker", err)
		return false, action.Response{Log: errors.Wrap(err, "unable to save the tracker").Error()}
	}
	ctx.Logger.Info("Vote added |  Validator : ", f.ValidatorAddress, " | Process Type : ", trackerlib.GetProcessTypeString(tracker.Type))
	yes, no := tracker.GetVotes()
	return true, action.Response{Log: "vote success, not ready to mint: " + strconv.Itoa(yes) + strconv.Itoa(no)}
}

func (reportFinalityMintTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.State.ConsumeVerifySigGas(1)
	ctx.State.ConsumeStorageGas(size)
	// check the used gas for the tx
	final := ctx.Balances.State.ConsumedGas()
	used := int64(final - start)
	ctx.Logger.Info("Gas Use : ",used)
	return true, action.Response{}
}

func mintTokens(ctx *action.Context, tracker *trackerlib.Tracker, oltTx ReportFinality) error {
	curr, ok := ctx.Currencies.GetCurrencyByName("ETH")
	if !ok {
		return errors.New("ETH currency not allowed")
	}
	lockAmount, err := ethereum.ParseLock(tracker.SignedETHTx)
	if err != nil {
		return err
	}

	oEthCoin := curr.NewCoinFromAmount(*balance.NewAmountFromBigInt(lockAmount.Amount))
	err = ctx.Balances.AddToAddress(oltTx.Locker, oEthCoin)
	if err != nil {
		ctx.Logger.Error(err)
		return errors.New("Unable to mint")
	}

	ethSupply := keys.Address(lockBalanceAddress)
	err = ctx.Balances.AddToAddress(ethSupply, oEthCoin)
	if err != nil {
		return errors.Wrap(err, "Unable to update total Eth supply")
	}

	tracker.State = trackerlib.Released
	err = ctx.ETHTrackers.Set(tracker)
	if err != nil {
		return err
	}
	return nil
}

func burnTokens(ctx *action.Context, tracker *trackerlib.Tracker, oltTx ReportFinality) error {

	tracker.State = trackerlib.Released
	err := ctx.ETHTrackers.Set(tracker)
	if err != nil {
		return err
	}

	return nil
}

func burnERC20Tokens(ctx *action.Context, tracker *trackerlib.Tracker, oltTx ReportFinality) error {
	ethTx, err := ethereum.DecodeTransaction(tracker.SignedETHTx)
	if err != nil {
		return err
	}
	ethOptions := ctx.ETHTrackers.GetOption()
	token, err := ethereum.GetToken(ethOptions.TokenList, *ethTx.To())
	if err != nil {
		return err
	}
	ctx.Logger.Info("Burn complete for token : " , token.TokName)
	tracker.State = trackerlib.Released
	err = ctx.ETHTrackers.Set(tracker)
	if err != nil {
		return err
	}

	return nil
}

func mintERC20tokens(ctx *action.Context, tracker *trackerlib.Tracker, oltTx ReportFinality) error {

	ethTx, err := ethereum.DecodeTransaction(tracker.SignedETHTx)
	if err != nil {
		return err
	}
	ethOptions := ctx.ETHTrackers.GetOption()
	token, err := ethereum.GetToken(ethOptions.TokenList, *ethTx.To())
	if err != nil {
		return err
	}
	ctx.Logger.Info("Minting Tokens of type : ", token.TokName)
	curr, ok := ctx.Currencies.GetCurrencyByName(token.TokName)
	if !ok {
		return errors.New("Currency not allowed ")
	}
	erc20Params,err := ethereum.ParseErc20Lock(ethOptions.TokenList,tracker.SignedETHTx)
	if err !=nil{
		return err
	}
	otokenCoin := curr.NewCoinFromAmount(*balance.NewAmountFromBigInt(erc20Params.TokenAmount))
	err = ctx.Balances.AddToAddress(oltTx.Locker, otokenCoin)
	if err != nil {
		ctx.Logger.Error(err)
		return errors.Errorf("Unable to mint token for : %s", token.TokName)
	}

	tokenSupply := keys.Address(TTClockBalanceAddress)
	err = ctx.Balances.AddToAddress(tokenSupply, otokenCoin)
	if err != nil{
		return errors.Errorf("Unable to update totalSupply for token : %s",token.TokName)
	}

	tracker.State = trackerlib.Released
	err = ctx.ETHTrackers.Set(tracker)
	if err != nil {
		return err
	}
	return nil
}
