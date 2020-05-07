/*
	Copyright 2017-2018 OneLedger

	Cover over the arguments of client requests
*/

// This file defines the functions available.
// Each of these should be able to be marshaled to and from JavaScript

package client

import (
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/tendermint/tendermint/libs/bytes"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
)

/*

We should divide our RPC layer using distinct services, which we can then define permissions & access rules over

Right now we only support request-reply styled responses.

Services:
	* broadcast
	* node
	* owner
	* query`json:"account"`
	* tx
*/

/* Blockchain service  */
type BalanceRequest struct {
	Address keys.Address `json:"address"`
}
type BalanceReply struct {
	// The balance of the account. Returns an empty balance
	// if the account is not found
	Balance string `json:"balance"`
	// The height when this balance was recorded
	Height int64 `json:"height"`
}

type ValidatorStatusRequest struct {
	Address keys.Address `json:"address"`
}

type ValidatorStatusReply struct {
	Height                int64  `json:"height"`
	Power                 int64  `json:"power"`
	Staking               string `json:"staking"`
	TotalDelegationAmount string `json:"totalDelegationAmount"`
	SelfDelegationAmount  string `json:"selfDelegationAmount"`
	DelegationAmount      string `json:"delegationAmount"`
	Exists                bool   `json:"exists"`
}

type DelegationStatusRequest struct {
	Address keys.Address `json:"address"`
}

type DelegationStatusReply struct {
	Balance                   string                   `json:"balance"`
	EffectiveDelegationAmount string                   `json:"effectiveDelegationAmount"`
	WithdrawableAmount        string                   `json:"withdrawableAmount"`
	MaturedAmounts            []*delegation.MatureData `json:"maturedAmount"`
}

/* Tx Service */

type SendTxRequest struct {
	From     keys.Address  `json:"from"`
	To       keys.Address  `json:"to"`
	Amount   action.Amount `json:"amount"`
	GasPrice action.Amount `json:"gasPrice"`
	Gas      int64         `json:"gas"`
}

type CreateTxReply struct {
	RawTx []byte `json:"rawTx"`
}

type StakeRequest struct {
	Address      keys.Address   `json:"address"`
	Amount       balance.Amount `json:"amount"`
	Name         string         `json:"name"`
	TmPubKeyType string         `json:"tmPubKeyType"`
	TmPubKey     []byte         `json:"tmPubKey"`
}

type StakeReply struct {
	RawTx []byte `json:"rawTx"`
}

type UnstakeRequest struct {
	Address keys.Address   `json:"address"`
	Amount  balance.Amount `json:"amount"`
}

type UnstakeReply struct {
	RawTx []byte `json:"rawTx"`
}

type WithdrawRequest struct {
	Address keys.Address   `json:"address"`
	Amount  balance.Amount `json:"amount"`
}

type WithdrawReply struct {
	RawTx []byte `json:"rawTx"`
}

type NodeNameRequest struct{}
type NodeNameReply struct {
	Name string `json:"name"`
}

type NodeAddressRequest struct{}
type NodeAddressReply struct {
	Address keys.Address `json:"address"`
}

type NodeIDRequest struct {
	ShouldShowIP bool `json:"shouldShowIP,omitempty"`
}
type NodeIDReply struct {
	Id string `json:"id"`
}

type AddAccountRequest = accounts.Account

type AddAccountReply struct {
	Account accounts.Account `json:"account"`
}

type GenerateAccountRequest struct {
	Name string `json:"name"`
}

type DeleteAccountRequest struct {
	Address keys.Address `json:"address"`
}
type DeleteAccountReply struct {
	Deleted bool `json:"deleted"`
}

type ListAccountsRequest struct{}
type ListAccountsReply struct {
	Accounts []accounts.Account `json:"accounts"`
}

type ListAccountAddressesReply struct {
	Addresses []string `json:"addresses"`
}

type ListValidatorsRequest struct{}
type ListValidatorsReply struct {
	// The list of active validators
	Validators []identity.Validator `json:"validators"`
	// Height at which this validator set was active
	Height int64           `json:"height"`
	VMap   map[string]bool `json:"vmap"`
}

type ListWitnessesRequest struct {
	ChainType chain.Type `json:"chainType"`
}
type ListWitnessesReply struct {
	// The list of active witnesses
	Witnesses []keys.Address `json:"witnesses"`
	// Height at which this witness set was active
	Height int64 `json:"height"`
}

type ListCurrenciesRequest struct{}
type ListCurrenciesReply struct {
	Currencies balance.Currencies `json:"currencies"`
}

type BroadcastRequest struct {
	RawTx     []byte         `json:"rawTx"`
	Signature []byte         `json:"signature"`
	PublicKey keys.PublicKey `json:"publicKey"`
}

type BroadcastReply struct {
	TxHash bytes.HexBytes `json:"txHash"`
	// OK indicates whether this broadcast was a request.
	// For TxSync, it indicates success of CheckTx. Does not guarantee inclusion of a block
	// For TxAsync, it always returns true
	// For TxCommit, it indicates the success of both CheckTx and DeliverTx. If the broadcast fails is false.
	OK     bool   `json:"ok"`
	Height *int64 `json:"height,omitempty"`
	Log    string `json:"log"`
}

type BroadcastMtSigRequest struct {
	RawTx      []byte             `json:"rawTx"`
	Signatures []action.Signature `json:"signatures"`
}

func (reply *BroadcastReply) FromResultBroadcastTx(result *ctypes.ResultBroadcastTx) {
	reply.TxHash = result.Hash
	reply.OK = result.Code == 0
	reply.Height = nil
	reply.Log = result.Log
}

func (reply *BroadcastReply) FromResultBroadcastTxCommit(result *ctypes.ResultBroadcastTxCommit) {
	reply.TxHash = result.Hash
	reply.OK = result.CheckTx.Code == 0 && result.DeliverTx.Code == 0
	reply.Height = &result.Height
	if len(result.CheckTx.Log) > 0 {
		reply.Log = result.CheckTx.Log + "[check]"
	} else if len(result.DeliverTx.Log) > 0 {
		reply.Log = result.DeliverTx.Log + "[deliver]"
	}
}

type NewAccountRequest struct {
	Name string `json:"name"`
}
type NewAccountReply struct {
	Account accounts.Account `json:"account"`
}

type SignRawTxRequest struct {
	RawTx   []byte         `json:"rawTx"`
	Address action.Address `json:"address"`
}

type SignRawTxResponse struct {
	Signature action.Signature `json:"signature"`
}

type BTCLockRequest struct {
	Txn         []byte        `json:"txn"`
	Address     keys.Address  `json:"address"`
	TrackerName string        `json:"trackerName"`
	GasPrice    action.Amount `json:"gasPrice"`
	Gas         int64         `json:"gas"`
}

type InputTransaction struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}
type BTCLockPrepareRequest struct {
	Inputs           []InputTransaction `json:"inputs"`
	AmountSatoshi    int64              `json:"amount"`
	FeeRate          int64              `json:"fee_rate"`
	ReturnAddressStr string             `json:"return_address"`
}

type BTCLockPrepareResponse struct {
	Txn         string `json:"txn"`
	TrackerName string `json:"trackerName"`
}

type BTCGetTrackerRequest struct {
	Name string `json:"name"`
}
type BTCGetTrackerReply struct {
	TrackerData string `json:"tracker"`
}

type BTCRedeemRequest struct {
	Address    keys.Address  `json:"address"`
	BTCAddress string        `json:"addressBTC"`
	Amount     int64         `json:"amount"`
	FeesBTC    int64         `json:"feesBTC"`
	GasPrice   action.Amount `json:"gasPrice"`
	Gas        int64         `json:"gas"`
}
type BTCRedeemPrepareResponse struct {
	RawTx       []byte `json:"rawTx"`
	TrackerName string `json:"trackerName"`
}

type ETHLockRequest struct {
	Txn     []byte
	Address keys.Address
	Fee     action.Amount `json:"fee"`
	Gas     int64         `json:"gas"`
}

type CurrencyBalanceRequest struct {
	Currency string       `json:"currency"`
	Address  keys.Address `json:"address"`
}
type CurrencyBalanceReply struct {
	Currency string `json:"currency"`
	Balance  string `json:"balance"`
	// The height when this balance was recorded
	Height int64 `json:"height"`
}

type EmptyRequest struct {
}

type MaxTrackerBalanceReply struct {
	MaxBalance int64 `json:"max_balance"`
}

type FeeOptionsReply struct {
	FeeOption fees.FeeOption `json:"feeOption"`
}

type ListTxTypesRequest struct{}
type ListTxTypesReply struct {
	TxTypes []action.TxTypeDescribe `json:"txTypes"`
}

type TxRequest struct {
	Hash  string `json:"hash"`
	Prove bool   `json:"prove"`
}

type TxResponse struct {
	Result ctypes.ResultTx `json:"result"`
}
