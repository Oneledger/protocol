/*
	Copyright 2017-2018 OneLedger

	Cover over the arguments of client requests
*/

// This file defines the functions available.
// Each of these should be able to be marshaled to and from JavaScript

package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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
type BalanceRequest = keys.Address
type BalanceReply struct {
	// The balance of the account. Returns an empty balance
	// if the account is not found
	Balance balance.Balance `json:"balance"`
	// The height when this balance was recorded
	Height int64 `json:"height"`
}

/* Tx Service */

type SendTxRequest struct {
	From   keys.Address  `json:"from"`
	To     keys.Address  `json:"to"`
	Amount action.Amount `json:"amount"`
	Fee    action.Amount `json:"fee"`
	Gas    int64         `json:"gas"`
}

type SendTxReply struct {
	RawTx []byte `json:"rawTx"`
}

type ApplyValidatorRequest struct {
	Address      keys.Address `json:"address"`
	Name         string       `json:"name"`
	Amount       string       `json:"amount"`
	Purge        bool         `json:"purge"`
	TmPubKeyType string       `json:"tmPubKeyType"`
	TmPubKey     []byte       `json:"tmPubKey"`
}

type ApplyValidatorReply struct {
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

type DeleteAccountRequest struct {
	Address keys.Address `json:"address"`
}
type DeleteAccountReply = bool

type ListAccountsRequest struct{}
type ListAccountsReply struct {
	Accounts []accounts.Account `json:"accounts"`
}

type ListValidatorsRequest struct{}
type ListValidatorsReply struct {
	// The list of active validators
	Validators []identity.Validator `json:"validators"`
	// Height at which this validator set was active
	Height int64 `json:"height"`
}

type BroadcastRequest struct {
	RawTx     []byte         `json:"rawTx"`
	Signature []byte         `json:"signature"`
	PublicKey keys.PublicKey `json:"publicKey"`
}

type BroadcastTxReply struct {
	TxHash keys.Address             `json:"txHash"`
	Result ctypes.ResultBroadcastTx `json:"result"`
}

type BroadcastTxCommitReply struct {
	TxHash keys.Address                   `json:"txHash"`
	Result ctypes.ResultBroadcastTxCommit `json:"result"`
}


/*
	ONS Request Types
*/
type ONSCreateRequest struct {
	Owner   keys.Address  `json:"owner"`
	Account keys.Address  `json:"account"`
	Name    string        `json:"name"`
	Price   action.Amount `json:"price"`
	Fee     action.Amount `json:"fee"`
	Gas     int64         `json:"gas"`
}

type ONSUpdateRequest struct {
	Owner   keys.Address  `json:"owner"`
	Account keys.Address  `json:"account"`
	Name    string        `json:"name"`
	Active  bool          `json:"active"`
	Fee     action.Amount `json:"fee"`
	Gas     int64         `json:"gas"`
}

type ONSSaleRequest struct {
	Name         string        `json:"name"`
	OwnerAddress keys.Address  `json:"owner"`
	Price        action.Amount `json:"price"`
	CancelSale   bool          `json:"cancel_sale"`
	Fee          action.Amount `json:"fee"`
	Gas          int64         `json:"gas"`
}

type ONSPurchaseRequest struct {
	Name     string        `json:"name"`
	Buyer    keys.Address  `json:"buyer"`
	Account  keys.Address  `json:"account"`
	Offering action.Amount `json:"offering"`
	Fee      action.Amount `json:"fee"`
	Gas      int64         `json:"gas"`
}

type ONSSendRequest struct {
	From   keys.Address  `json:"from"`
	Name   string        `json:"name"`
	Amount action.Amount `json:"amount"`
	Fee    action.Amount `json:"fee"`
	Gas    int64         `json:"gas"`

type SignRawTxRequest struct {
	RawTx   []byte         `json:"rawTx"`
	Address action.Address `json:"address"`
}

type SignRawTxResponse struct {
	Signature action.Signature `json:"signature"`
}
