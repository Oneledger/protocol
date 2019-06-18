/*
	Copyright 2017-2018 OneLedger

	Cover over the arguments of client requests
*/

// Each of these should be able to be marshaled to and from javascript

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

Right now we only support request-reply styled responses

Services:
	* NodeService - Node-related queries and requests
	* BlockchainService - General-blockchain
*/

/* Blockchain service  */
type BalanceRequest = keys.Address
type BalanceReply struct {
	Balance balance.Balance `json:"balance"`
	// The height perfect
	Height int64 `json:"height"`
}

type SendTxRequest struct {
	From   keys.Address  `json:"from"`
	To     keys.Address  `json:"to"`
	Amount action.Amount `json:"amount"`
	// Unused
	Fee action.Amount `json:"fee"`
	Gas int64         `json:"gas"`
}
type SendTxReply struct {
	RawTx []byte `json:"rawTx"`
}

// SendRawTx is for broadcasting a raw self-signed transaction over the network
type SendRawTxRequest struct {
	// Msg is the raw transaction bytes
	RawTx     []byte         `json:"rawTx"`
	Signature []byte         `json:"signature"`
	PublicKey keys.PublicKey `json:"publicKey"`
}
type SendRawTxReply struct {
	// The result of broadcasting this transaction to the network
	Result ctypes.ResultBroadcastTx `json:"result"`
}

/* These are node-related access requests */
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

type ApplyValidatorRequest struct {
	Address      keys.Address `json:"address"`
	Name         string       `json:"name"`
	Amount       string       `json:"amount"`
	Purge        bool         `json:"purge"`
	TmPubKeyType string       `json:"tmPubKeyType"`
	TmPubKey     []byte       `json:"tmPubKey"`
}
type ApplyValidatorReply struct{}

type ListValidatorsRequest struct{}
type ListValidatorsReply struct {
	// The list of active validators
	Validators []identity.Validator `json:"validators"`
	// Height at which this validator set was active
	Height int64 `json:"height"`
}
