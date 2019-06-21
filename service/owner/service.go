package owner

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/rpc"
	"github.com/pkg/errors"
)

type Service struct {
	accounts accounts.Wallet
	logger   *log.Logger
}

func NewService(accts accounts.Wallet, logger *log.Logger) *Service {
	return &Service{
		accounts: accts,
	}
}

func Name() string {
	return "owner"
}

// AddAccount adds an account to the local accounts store of the node
func (svc *Service) AddAccount(acc client.AddAccountRequest, reply *client.AddAccountReply) error {
	err := svc.accounts.Add(acc)
	if err != nil {
		err := errors.Wrap(err, "error in adding account to walletstore")
		return rpc.InternalError(err.Error())
	}

	acct, err := svc.accounts.GetAccount(acc.Address())
	*reply = client.AddAccountReply{Account: acct}
	return nil
}

// DeleteAccount deletes an account from the local store
func (svc *Service) DeleteAccount(req client.DeleteAccountRequest, reply *client.DeleteAccountReply) error {
	var nilAccount accounts.Account
	toDelete, err := svc.accounts.GetAccount(req.Address)
	if err != nil || toDelete == nilAccount {
		return rpc.NewError(rpc.CodeNotFound, "account doesn't exist")
	}
	err = svc.accounts.Delete(toDelete)
	if err != nil {
		return rpc.InternalError("error in deleting account from walletstore")
	}

	*reply = true
	return nil
}

// ListAccounts lists all accounts available in the local store
func (svc *Service) ListAccounts(req client.ListAccountsRequest, reply *client.ListAccountsReply) error {
	accts := svc.accounts.Accounts()
	if accts == nil {
		accts = make([]accounts.Account, 0)
	}
	*reply = client.ListAccountsReply{Accounts: accts}

	return nil
}

// ListAccountAddresses lists all accounts available in the local store
func (svc *Service) ListAccountAddresses(req client.ListAccountsRequest, reply *client.ListAccountAddressesReply) error {
	accts := svc.accounts.Accounts()
	if accts == nil {
		accts = make([]accounts.Account, 0)
	}
	addrs := make([]string, len(accts))
	for i := range accts {
		addrs[i] = accts[i].Address().Humanize()
	}
	*reply = client.ListAccountAddressesReply{Addresses: addrs}

	return nil
}

func (svc *Service) SignWithAddress(req client.SignRawTxRequest, reply *client.SignRawTxResponse) error {
	pkey, signed, err := svc.accounts.SignWithAddress(req.RawTx, req.Address)
	if err != nil {
		return rpc.InternalError(err.Error())
	}
	*reply = client.SignRawTxResponse{Signature: action.Signature{Signed: signed, Signer: pkey}}
	return nil
}

func (svc *Service) NewAccount(req client.NewAccountRequest, reply *client.NewAccountReply) error {

	pubKey, privKey, err := keys.NewKeyPairFromTendermint()
	if err != nil {
		return rpc.InternalError(err.Error())
	}

	ct, err := chain.TypeFromName("OneLedger")
	if err != nil {
		return rpc.InternalError(err.Error())
	}

	acc, err := accounts.NewAccount(ct, req.Name, &privKey, &pubKey)
	if err != nil {
		return rpc.InternalError(err.Error())
	}

	reply = &client.NewAccountReply{acc}
	return nil
}
