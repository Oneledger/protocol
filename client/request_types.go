/*
	Copyright 2017-2018 OneLedger

	Cover over the arguments of client requests
*/

// This file defines the functions available.
// Each of these should be able to be marshaled to and from JavaScript

package client

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
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
	* query
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

/* Tx Service */

type SendTxRequest struct {
	From     keys.Address  `json:"from"`
	To       keys.Address  `json:"to"`
	Amount   action.Amount `json:"amount"`
	GasPrice action.Amount `json:"gasprice"`
	Gas      int64         `json:"gas"`
}

type SendTxReply struct {
	RawTx []byte `json:"rawTx"`
}

type ApplyValidatorRequest struct {
	Address      keys.Address   `json:"address"`
	Name         string         `json:"name"`
	Amount       balance.Amount `json:"amount"`
	Purge        bool           `json:"purge"`
	TmPubKeyType string         `json:"tmPubKeyType"`
	TmPubKey     []byte         `json:"tmPubKey"`
}

type ApplyValidatorReply struct {
	RawTx []byte `json:"rawTx"`
}

type WithdrawRewardRequest struct {
	From     keys.Address  `json:"from"`
	To       keys.Address  `json:"to"`
	GasPrice action.Amount `json:"gasprice"`
	Gas      int64         `json:"gas"`
}

type WithdrawRewardReply struct {
	RawTx []byte `json:"rawTx"`
}

type NodeNameRequest struct{}
type NodeNameReply = string

type NodeAddressRequest struct{}
type NodeAddressReply = keys.Address

type NodeIDRequest struct {
	ShouldShowIP bool `json:"shouldShowIP,omitempty"`
}
type NodeIDReply = string

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
type DeleteAccountReply = bool

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
	TxHash action.Address `json:"txHash"`
	// OK indicates whether this broadcast was a request.
	// For TxSync, it indicates success of CheckTx. Does not guarantee inclusion of a block
	// For TxAsync, it always returns true
	// For TxCommit, it indicates the success of both CheckTx and DeliverTx. If the broadcast fails is false.
	OK     bool   `json:"ok"`
	Height *int64 `json:"height,omitempty"`
	Log    string `json:"log"`
}

func (reply *BroadcastReply) FromResultBroadcastTx(result *ctypes.ResultBroadcastTx) {
	reply.TxHash = action.Address(result.Hash)
	reply.OK = result.Code == 0
	reply.Height = nil
	reply.Log = result.Log
}

func (reply *BroadcastReply) FromResultBroadcastTxCommit(result *ctypes.ResultBroadcastTxCommit) {
	reply.TxHash = action.Address(result.Hash)
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
	Signature   []byte        `json:"signature"`
	Address     keys.Address  `json:"address"`
	TrackerName string        `json:"tracker_name"`
	GasPrice    action.Amount `json:"gasprice"`
	Gas         int64         `json:"gas"`
}

type BTCLockPrepareRequest struct {
	Hash    string `json:"hash"`
	Index   uint32 `json:"index"`
	FeesBTC int64  `json:"fees_btc"`
}
type BTCLockPrepareResponse struct {
	Txn         string `json:"txn"`
	TrackerName string `json:"tracker_name"`
}

type BTCGetTrackerRequest struct {
	Name string `json:"name"`
}
type BTCGetTrackerReply struct {
	Tracker bitcoin.Tracker `json:"tracker"`
}

type BTCLockRedeemRequest struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
	FeesBTC int64  `json:"fees_btc"`
}
type BTCRedeemPrepareResponse struct {
	Txn         string `json:"txn"`
	TrackerName string `json:"tracker_name"`
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
