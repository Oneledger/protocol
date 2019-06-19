package client

import (
	"errors"

	"github.com/Oneledger/protocol/rpc"
)

// A generic client for accessing rpc services.
// Eventually each service will be broken out onto its own type
// TODO: The methods defined here should handle context.Context
type ServiceClient struct {
	*rpc.Client
}

/*
	Blockchain Service
	- Query basic information from the blockchain.
	- Explorer stuff. Get tx by hash, query acct balance, get a block by height or list of blocks.
*/

func (c *ServiceClient) Balance(request BalanceRequest) (out BalanceReply, err error) {
	if len(request) <= 20 {
		return out, errors.New("address has insufficient length")
	}
	err = c.Call("server.Balance", request, &out)
	return
}

/*
	Node Service
	- Query information from the chain from us
*/

func (c *ServiceClient) NodeName() (out NodeNameReply, err error) {
	err = c.Call("server.NodeName", struct{}{}, &out)
	return
}

func (c *ServiceClient) NodeAddress() (out NodeAddressReply, err error) {
	err = c.Call("server.NodeAddress", struct{}{}, &out)
	return
}

func (c *ServiceClient) NodeID(req NodeIDRequest) (out NodeIDReply, err error) {
	err = c.Call("server.NodeID", req, &out)
	return
}

/*
	Tx Service
	- Returns raw messages to be signed so you can broadcast transactions
	- TODO: There should be methods specific for having the node sign & broadcast for us in one go, but the authentication
	  part needs to be fleshed out.
*/

func (c *ServiceClient) SendTx(req SendTxRequest) (out SendTxReply, err error) {
	err = c.Call("server.SendTx", req, &out)
	return
}

/*
	Accounts/Identities Service
	- Manages local accounts and what they deserve to be
*/

func (c *ServiceClient) RawTxBroadcast(req RawTxBroadcastRequest) (out RawTxBroadcastReply, err error) {
	err = c.Call("server.RawTxBroadcast", req, &out)
	return
}

func (c *ServiceClient) AddAccount(req AddAccountRequest) (out AddAccountReply, err error) {
	err = c.Call("server.AddAccount", req, &out)
	return
}

func (c *ServiceClient) DeleteAccount(req DeleteAccountRequest) (out DeleteAccountReply, err error) {
	err = c.Call("server.DeleteAccount", req, &out)
	return
}

func (c *ServiceClient) ListAccounts() (out ListAccountsReply, err error) {
	err = c.Call("server.ListAccounts", struct{}{}, &out)
	return
}

func (c *ServiceClient) ApplyValidator(req ApplyValidatorRequest) (out ApplyValidatorReply, err error) {
	err = c.Call("server.ApplyValidator", req, &out)
	return
}

func (c *ServiceClient) ListValidators() (out ListValidatorsReply, err error) {
	err = c.Call("server.ListValidators", struct{}{}, &out)
	return
}

/*
	Broadcast Service
	- Broadcast transactions to the network. This should be protected
*/
