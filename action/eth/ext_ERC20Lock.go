//Package for transactions related to Etheruem
package eth

import (
	"encoding/json"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	ethchaindriver "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/ethereum"
)

// Lock is a struct for one-Ledger transaction for ERC20 Lock
type ERC20Lock struct {
	Locker action.Address
	ETHTxn []byte // Raw Transaction for Locking Tokens
}

var _ action.Msg = &ERC20Lock{}
//Signers return the Address of the owner who created the transaction
func (E ERC20Lock) Signers() []action.Address {
	return []action.Address{E.Locker}
}

// Type returns the type of current action
func (E ERC20Lock) Type() action.Type {
	return action.ERC20_LOCK
}
// Tags creates the tags to associate with the transaction
func (E ERC20Lock) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(E.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.locker"),
		Value: E.Locker.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

//Marshal ERC20Lock to byte array
func (E ERC20Lock) Marshal() ([]byte, error) {
	return json.Marshal(E)
}

func (E *ERC20Lock) Unmarshal(data []byte) error {
	return json.Unmarshal(data, E)
}

type ethERC20LockTx struct {
}

var _ action.Tx = ethERC20LockTx{}
// Validate provides basic validation for transaction Type and Fee
func (e ethERC20LockTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	// unmarshal the tx message
	erclock := &ERC20Lock{}
	err := erclock.Unmarshal(signedTx.Data)
	if err != nil {
		ctx.Logger.Error("error unmarshalling Data field of ERC LOCK trasaction")
		return false, errors.Wrap(err, action.ErrWrongTxType.Msg)
	}
	err = action.ValidateBasic(signedTx.RawBytes(), erclock.Signers(), signedTx.Signatures)
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
	return true, nil
}
// ProcessCheck runs checks on the transaction without commiting it .
func (e ethERC20LockTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runERC20Lock(ctx, tx)
}
// ProcessDeliver run checks on transaction and commits it to a new block
func (e ethERC20LockTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runERC20Lock(ctx, tx)
}
// ProcessFee process the transaction Fee in OLT
func (e ethERC20LockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.State.ConsumeUpfront(237600)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
// runERC20Lock has the common functionality for ProcessCheck and ProcessDeliver
// Provides security checks for transaction
func runERC20Lock(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	erc20lock := &ERC20Lock{}
	err := erc20lock.Unmarshal(tx.Data)
	if err != nil {
		ctx.Logger.Error("wrong tx type", err)
		return false, action.Response{Log: "wrong tx type"}
	}
	ethTx, err := ethchaindriver.DecodeTransaction(erc20lock.ETHTxn)
	if err != nil {
		ctx.Logger.Error("decode eth txn err", err)
		return false, action.Response{
			Log: "decode eth txn error" + err.Error(),
		}
	}
	ethOptions := ctx.ETHTrackers.GetOption()
	token, err := ethchaindriver.GetToken(ethOptions.TokenList, *ethTx.To())
	if err != nil {
		return false, action.Response{
			Log: err.Error(),
		}
	}
	ok, err := ethchaindriver.VerfiyERC20Lock(erc20lock.ETHTxn, token.TokAbi, ethOptions.ERCContractAddress)
	if err != nil {
		ctx.Logger.Error("Unable to verify ERC LOCK transaction")
		return false, action.Response{
			Log: "Unable to verify transaction" + err.Error(),
		}
	}
	if !ok {
		ctx.Logger.Error("To field of Transaction does not match OneLedger Contract Address")
		return false, action.Response{
			Log: "To field of Transaction does not match OneLedger Contract Address" + err.Error(),
		}
	}
	validatorList, err := ctx.Validators.GetValidatorsAddress()
	if err != nil {
		ctx.Logger.Error("err in getting validator address", err)
		return false, action.Response{Log: "error in getting validator addresses" + err.Error()}
	}
	curr, ok := ctx.Currencies.GetCurrencyByName(token.TokName)
	if !ok {
		return false, action.Response{Log: fmt.Sprintf("Token not Supported : %s ", token.TokName)}
	}
	erc20Params,err := ethchaindriver.ParseErc20Lock(ethOptions.TokenList,erc20lock.ETHTxn)
	if err !=nil{
		return false,action.Response{
			Log: err.Error(),
		}
	}
	lockToken := curr.NewCoinFromString(erc20Params.TokenAmount.String())
	// Adding lock amount to common address to maintain count of total oToken minted
	tokenSupply := action.Address(TTClockBalanceAddress)
	balCoin, err := ctx.Balances.GetBalanceForCurr(tokenSupply, &curr)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("Unable to get Eth lock total balance %s", erc20lock.Locker)}
	}
	totalSupplyToken := curr.NewCoinFromString(totalTTCSupply)
	if !balCoin.Plus(lockToken).LessThanEqualCoin(totalSupplyToken) {
		return false, action.Response{Log: fmt.Sprintf("Token lock exceeded limit ,for Token : %s ", token.TokName)}
	}
	
	
	tracker := ethereum.NewTracker(
		ethereum.ProcessTypeLockERC,
		erc20lock.Locker,
		erc20lock.ETHTxn,
		ethcommon.BytesToHash(erc20lock.ETHTxn),
		validatorList,
	)

	err = ctx.ETHTrackers.Set(tracker)
	if err != nil {
		ctx.Logger.Error("error saving eth tracker", err)
		return false, action.Response{Log: "error saving eth tracker: " + err.Error()}
	}

	return true, action.Response{
		Tags: erc20lock.Tags(),
	}
}
