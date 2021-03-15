//Package for transactions related to Etheruem
package eth

import (
	"encoding/json"

	"github.com/tendermint/tendermint/libs/kv"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &Redeem{}

// Redeem is a struct for one-Ledger transaction for Ether Redeem
type Redeem struct {
	Owner  action.Address    //User Oneledger address
	To     ethcommon.Address //User Ethereum address
	ETHTxn []byte
}

//Signers return the Address of the owner who created the transaction
func (r Redeem) Signers() []action.Address {
	return []action.Address{r.Owner}
}

// Type returns the type of current action
func (r Redeem) Type() action.Type {
	return action.ETH_REDEEM
}

// Tags creates the tags to associate with the transaction
func (r Redeem) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(r.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: r.Owner,
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.tracker"),
		Value: ethcommon.BytesToHash(r.ETHTxn).Bytes(),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

//Marshal Redeem to byte array
func (r Redeem) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Redeem) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

var _ action.Tx = ethRedeemTx{}

type ethRedeemTx struct {
}

// Validate provides basic validation for transaction Type and Fee
func (ethRedeemTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	redeem := &Redeem{}
	err := redeem.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(err, action.ErrWrongTxType.Error())
	}
	err = action.ValidateBasic(signedTx.RawBytes(), redeem.Signers(), signedTx.Signatures)
	if err != nil {
		ctx.Logger.Error("validate basic failed", err)
		return false, err
	}

	// validate fee
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		ctx.Logger.Error("validate fee failed", err)
		return false, err
	}

	if redeem.ETHTxn == nil {
		ctx.Logger.Error("eth txn is nil")
		return false, action.ErrMissingData
	}

	return true, nil

}

// ProcessCheck runs checks on the transaction without commiting it .
func (ethRedeemTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Executing Processcheck for ETH REDEEM")
	return runRedeem(ctx, tx)

}

// ProcessDeliver run checks on transaction and commits it to a new block
func (ethRedeemTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Executing Processdeliver for ETH REDEEM")
	return runRedeem(ctx, tx)
}

// runRedeem has the common functionality for ProcessCheck and ProcessDeliver
// Provides security checks for transaction
func runRedeem(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	redeem := &Redeem{}
	err := redeem.Unmarshal(tx.Data)
	if err != nil {
		ctx.Logger.Error("")
		return false, action.Response{Log: errors.Wrap(action.ErrUnserializable, err.Error()).Error()}
	}
	ethOptions, err := ctx.GovernanceStore.GetETHChainDriverOption()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, gov.ErrGetEthOptions, redeem.Tags(), err)
	}
	req, err := ethereum.ParseRedeem(redeem.ETHTxn, ethOptions.ContractABI)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidExtTx, redeem.Tags(), err)
	}

	c, ok := ctx.Currencies.GetCurrencyByName("ETH")
	if !ok {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidCurrency, redeem.Tags(), err)
	}

	coin := c.NewCoinFromAmount(*balance.NewAmountFromBigInt(req.Amount))
	err = ctx.Balances.MinusFromAddress(redeem.Owner, coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorMinusFailed, redeem.Tags(), err)
	}
	// Subtracting from common address to maintain count of the total oEth minted
	ethSupply := keys.Address(ethOptions.TotalSupplyAddr)
	err = ctx.Balances.MinusFromAddress(ethSupply, coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorMinusFailed, redeem.Tags(), err)
	}

	witnesses, err := ctx.Witnesses.GetWitnessAddresses(chain.ETHEREUM)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrGettingWitnessList, redeem.Tags(), err)
	}
	name := ethcommon.BytesToHash(redeem.ETHTxn)
	if ctx.ETHTrackers.WithPrefixType(trackerlib.PrefixOngoing).Exists(name) || ctx.ETHTrackers.WithPrefixType(trackerlib.PrefixFailed).Exists(name) || ctx.ETHTrackers.WithPrefixType(trackerlib.PrefixPassed).Exists(name) {
		return helpers.LogAndReturnFalse(ctx.Logger, trackerlib.ErrETHTrackerExists, redeem.Tags(), errors.New("Tracker with same TXHASH already exists"))
	}

	tracker := trackerlib.NewTracker(
		trackerlib.ProcessTypeRedeem,
		redeem.Owner,
		redeem.ETHTxn,
		name,
		witnesses,
	)

	tracker.State = trackerlib.New
	tracker.ProcessOwner = redeem.Owner
	tracker.SignedETHTx = redeem.ETHTxn
	tracker.To = redeem.To.Bytes()

	// Save eth Tracker
	err = ctx.ETHTrackers.WithPrefixType(trackerlib.PrefixOngoing).Set(tracker)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, trackerlib.ErrETHTrackerUnableToSet, redeem.Tags(), err)
	}
	ctx.Logger.Debug("Redeem Tracker set | Jobs Starting now")
	return true, action.Response{
		Data:      nil,
		Log:       "",
		Info:      "Transaction received ,Redeem in progress",
		GasWanted: 0,
		GasUsed:   0,
		Events:    action.GetEvent(redeem.Tags(), "eth_redeem"),
	}
}

// ProcessFee process the transaction Fee in OLT
func (ethRedeemTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.State.ConsumeUpfront(250400)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
