//Package for transactions related to Etheruem
package eth

import (
	"encoding/json"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &ERC20Redeem{}

// Lock is a struct for one-Ledger transaction for ERC20 Redeem
type ERC20Redeem struct {
	Owner  action.Address //User Oneledger address
	To     action.Address //User Ethereum address
	ETHTxn []byte
}

//Signers return the Address of the owner who created the transaction
func (E ERC20Redeem) Signers() []action.Address {
	return []action.Address{E.Owner}
}

// Type returns the type of current action
func (E ERC20Redeem) Type() action.Type {
	return action.ERC20_REDEEM
}

// Tags creates the tags to associate with the transaction
func (E ERC20Redeem) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(E.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: E.Owner,
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.tracker"),
		Value: ethcommon.BytesToHash(E.ETHTxn).Bytes(),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

//Marshal ERC20Redeem to byte array
func (E ERC20Redeem) Marshal() ([]byte, error) {
	return json.Marshal(E)
}

func (E *ERC20Redeem) Unmarshal(data []byte) error {
	return json.Unmarshal(data, E)
}

var _ action.Tx = ethERC20RedeemTx{}

type ethERC20RedeemTx struct {
}

// Validate provides basic validation for transaction Type and Fee
func (e ethERC20RedeemTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	erc20redeem := &ERC20Redeem{}
	err := erc20redeem.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(err, action.ErrWrongTxType.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), erc20redeem.Signers(), signedTx.Signatures)
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

	if erc20redeem.ETHTxn == nil {
		ctx.Logger.Error("eth txn is nil")
		return false, action.ErrMissingData
	}
	return true, nil
}

// ProcessCheck runs checks on the transaction without commiting it .
func (e ethERC20RedeemTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runERC20Reddem(ctx, tx)
}

// ProcessDeliver run checks on transaction and commits it to a new block
func (e ethERC20RedeemTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runERC20Reddem(ctx, tx)
}

// ProcessFee process the transaction Fee in OLT
func (e ethERC20RedeemTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{Log: "ProcessFee"}
}

// runERC20Redeem has the common functionality for ProcessCheck and ProcessDeliver
// Provides security checks for transaction
func runERC20Reddem(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	erc20redeem := &ERC20Redeem{}
	err := erc20redeem.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: action.ErrUnserializable.Error()}
	}

	ethOptions := ctx.ETHTrackers.GetOption()
	redeemParams, err := ethereum.ParseERC20RedeemParams(erc20redeem.ETHTxn, ethOptions.ERCContractABI)
	if err != nil {
		ctx.Logger.Error(err)
		return false, action.Response{Log: action.ErrTokenNotSupported.Error()}
	}

	token, err := ethereum.ParseERC20RedeemToken(erc20redeem.ETHTxn, ethOptions.TokenList, ethOptions.ERCContractABI)
	if err != nil {
		ctx.Logger.Error(err)
		return false, action.Response{Log: action.ErrTokenNotSupported.Error()}
	}

	c, ok := ctx.Currencies.GetCurrencyByName(token.TokName)
	if !ok {
		return false, action.Response{Log: "Token not registered "}
	}

	coin := c.NewCoinFromAmount(*balance.NewAmountFromBigInt(redeemParams.Amount))
	err = ctx.Balances.MinusFromAddress(erc20redeem.Owner, coin)
	if err != nil {
		fmt.Println("Not enough funds")
		return false, action.Response{Log: action.ErrNotEnoughFund.Error()}
	}

	// Subtracting from common address to maintain count of the total oToken minted
	tokenSupply := keys.Address(lockBalanceAddress)
	err = ctx.Balances.MinusFromAddress(tokenSupply, coin)
	if err != nil {
		return false, action.Response{Log: action.ErrNotEnoughFund.Error()}
	}

	validators, err := ctx.Validators.GetValidatorsAddress()
	if err != nil {
		return false, action.Response{Log: "error in getting validator addresses" + err.Error()}
	}
	name := ethcommon.BytesToHash(erc20redeem.ETHTxn)
	if ctx.ETHTrackers.Exists(name) {
		return false, action.Response{
			Log: "Tracker already exists",
		}
	}

	tracker := trackerlib.NewTracker(
		trackerlib.ProcessTypeRedeemERC,
		erc20redeem.Owner,
		erc20redeem.ETHTxn,
		name,
		validators,
	)

	tracker.State = trackerlib.New
	tracker.ProcessOwner = erc20redeem.Owner
	tracker.SignedETHTx = erc20redeem.ETHTxn
	tracker.To = erc20redeem.To

	// Save eth Tracker
	err = ctx.ETHTrackers.Set(tracker)
	return true, action.Response{
		Data:      nil,
		Log:       "",
		Info:      "Transaction received ,Redeem in progress",
		GasWanted: 0,
		GasUsed:   0,
		Tags:      erc20redeem.Tags(),
	}
}
