/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package app

import (
	"fmt"
	"io"
	"net/url"
	"runtime/debug"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/p2p"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type RPCServerContext struct {
	nodeName     string
	balances     *balance.Store
	accounts     accounts.Wallet
	currencies   *balance.CurrencyList
	cfg          config.Server
	nodeContext  NodeContext
	validatorSet *identity.ValidatorStore

	services *client.ServiceContext
	logger   *log.Logger
}

func NewClientHandler(
	nodeName string,
	balances *balance.Store,
	accounts accounts.Wallet,
	currencies *balance.CurrencyList,
	cfg config.Server,
	nodeContext NodeContext,
	validatorSet *identity.ValidatorStore,
	logWriter io.Writer,
) *RPCServerContext {
	return &RPCServerContext{
		nodeName:     nodeName,
		balances:     balances,
		accounts:     accounts,
		currencies:   currencies,
		cfg:          cfg,
		nodeContext:  nodeContext,
		validatorSet: validatorSet,
		logger:       log.NewLoggerWithPrefix(logWriter, "client_Handler"),
	}
}

// NodeName returns the name of a node. This is useful for displaying it at cmdline.
func (h *RPCServerContext) NodeName(_ client.NodeNameRequest, reply *client.NodeNameReply) error {
	*reply = h.nodeName
	return nil
}

func (h *RPCServerContext) NodeAddress(_ client.NodeAddressRequest, reply *client.NodeAddressReply) error {
	*reply = h.nodeContext.Address()
	return nil
}

func (h *RPCServerContext) NodeID(req client.NodeIDRequest, reply *client.NodeIDReply) error {
	nodeKey, err := p2p.LoadNodeKey(h.cfg.TMConfig().NodeKeyFile())
	if err != nil {
		return errors.Wrap(err, "error loading node key")
	}

	// silenced error because not present means false
	ip := p2pAddressFromCFG(h.cfg)
	if req.ShouldShowIP {
		u, err := url.Parse(ip)
		if err != nil {
			return errors.Wrap(err, "error in parsing configured url")
		}
		out := fmt.Sprintf("%s@%s", nodeKey.ID(), u.Host)
		*reply = out
	} else {
		*reply = string(nodeKey.ID())
	}
	return nil
}

// This function returns the external p2p address if it exists, but falls back to the regular p2p address if it is
// not present from the config
func p2pAddressFromCFG(cfg config.Server) string {
	extP2P := cfg.Network.ExternalP2PAddress
	if extP2P != "" {
		return cfg.Network.P2PAddress
	}

	return cfg.Network.ExternalP2PAddress
}

// GetBalance gets the balance of an address
// TODO make it more generic to handle account name and identity
func (h *RPCServerContext) Balance(req client.BalanceRequest, resp *client.BalanceReply) error {
	addr := req
	bal, err := h.balances.Get(addr, true)

	if err != nil && err == balance.ErrNoBalanceFoundForThisAddress {
		bal = balance.NewBalance()
	} else if err != nil {
		h.logger.Error("error getting balance", err)
		return errors.Wrap(err, "error getting balance")
	}

	*resp = client.BalanceReply{
		Balance: *bal,
		Height:  h.balances.Version,
	}
	return nil
}

// AddAccount adds an account to accounts store of the node
func (h *RPCServerContext) AddAccount(acc client.AddAccountRequest, reply *client.AddAccountReply) error {
	err := h.accounts.Add(acc)
	if err != nil {
		return errors.Wrap(err, "error in adding account to walletstore")
	}

	acct, err := h.accounts.GetAccount(acc.Address())
	*reply = client.AddAccountReply{Account: acct}
	return nil
}

// DeleteAccount deletes an account from the accounts store of node
func (h *RPCServerContext) DeleteAccount(req client.DeleteAccountRequest, reply *client.DeleteAccountReply) error {
	defer h.recoverPanic()
	// TODO: Need to verify that we're allowed to delete this account.
	//       The input should be a pair of both the account address and a signed version of the address.
	//       Only delete the account if the signature is valid
	var nilAccount accounts.Account
	toDelete, err := h.accounts.GetAccount(req.Address)
	if err != nil || toDelete == nilAccount {
		return errors.New("account doesn't exist!")
	}
	err = h.accounts.Delete(toDelete)
	if err != nil {
		return errors.Wrap(err, "error in deleting account from walletstore")
	}

	*reply = true
	return nil
}

// ListAccounts returns a list of all accounts in the accounts store of node
func (h *RPCServerContext) ListAccounts(req client.ListAccountsRequest, reply *client.ListAccountsReply) error {
	// TODO: pagination

	accts := h.accounts.Accounts()
	if accts == nil {
		accts = make([]accounts.Account, 0)
	}
	*reply = client.ListAccountsReply{Accounts: accts}

	return nil
}

func (h *RPCServerContext) SendTx(args client.SendTxRequest, reply *client.SendTxReply) error {
	send := action.Send{
		From:   keys.Address(args.From),
		To:     keys.Address(args.To),
		Amount: args.Amount,
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.BaseTx{
		Data: send,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	if _, err := h.accounts.GetAccount(args.From); err != nil {
		return errors.New("Account doesn't exist. Send a raw tx instead")
	}

	pubKey, signed, err := h.accounts.SignWithAddress(tx.Bytes(), send.From)
	if err != nil {
		return err
	}
	tx.Signatures = []action.Signature{{pubKey, signed}}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return errors.Wrap(err, "err while network serialization")
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

// SendRawTx.
func (h *RPCServerContext) RawTxBroadcast(args client.RawTxBroadcastRequest, reply *client.RawTxBroadcastReply) error {
	var act action.BaseTx

	signer, err := args.PublicKey.GetHandler()
	if err != nil {
		return errors.New("not a valid public key")
	}

	isVerified := signer.VerifyBytes(args.RawTx, args.Signature)
	if !isVerified {
		return errors.New("signature is not valid")
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(act)
	if err != nil {
		return errors.New("failed to serialize transaction")
	}

	result, err := h.broadcastTxSync(packet)
	if err != nil {
		return errors.Wrap(err, "failed to broadcast")
	}

	*reply = client.RawTxBroadcastReply{
		Result: *result,
	}
	return nil
}

func (h *RPCServerContext) ApplyValidator(args client.ApplyValidatorRequest, reply *client.ApplyValidatorReply) error {
	defer h.recoverPanic()

	if len(args.Name) < 1 {
		args.Name = h.nodeName
	}

	if len(args.Address) < 1 {
		handler, err := h.nodeContext.PubKey().GetHandler()
		if err != nil {
			return err
		}
		args.Address = handler.Address()
	}

	pubkey := &keys.PublicKey{keys.GetAlgorithmFromTmKeyName(args.TmPubKeyType), args.TmPubKey}
	if len(args.TmPubKey) < 1 {
		*pubkey = h.nodeContext.ValidatorPubKey()
	}

	handler, err := pubkey.GetHandler()
	if err != nil {

		return err
	}

	addr := handler.Address()
	apply := action.ApplyValidator{
		Address:          keys.Address(args.Address),
		Stake:            action.Amount{Currency: "VT", Value: args.Amount},
		NodeName:         args.Name,
		ValidatorAddress: addr,
		ValidatorPubKey:  *pubkey,
	}

	uuidNew, _ := uuid.NewUUID()
	tx := action.BaseTx{
		Data: apply,
		Fee:  action.Fee{action.Amount{Currency: "OLT", Value: "0.1"}, 1},
		Memo: uuidNew.String(),
	}

	pubKey, signed, err := h.accounts.SignWithAccountIndex(tx.Bytes(), 0)
	if err != nil {
		return err
	}
	tx.Signatures = []action.Signature{{pubKey, signed}}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return errors.Wrap(err, "err while network serialization")
	}

	*reply = client.ApplyValidatorReply{packet}

	return nil
}

// ListValidator returns a list of all validator
func (h *RPCServerContext) ListValidators(_ client.ListValidatorsRequest, reply *client.ListValidatorsReply) error {
	validators, err := h.validatorSet.GetValidatorSet()
	if err != nil {
		return errors.Wrap(err, "err while retrieving validators info")
	}

	*reply = client.ListValidatorsReply{
		Validators: validators,
		Height:     h.balances.Version,
	}

	return nil
}

/*
	Client handler util methods
*/
// recoverPanic common method to recover from any handler panic
func (h *RPCServerContext) recoverPanic() {
	if r := recover(); r != nil {
		h.logger.Error("recovering a panic")
		debug.PrintStack()
	}
}

// broadcastTxAsync asynchronously broadcasts a transaction to Tendermint
func (h *RPCServerContext) broadcastTxAsync(packet []byte) (res *ctypes.ResultBroadcastTx, err error) {
	return h.services.BroadcastTxAsync(packet)
}

func (h *RPCServerContext) broadcastTxSync(packet []byte) (res *ctypes.ResultBroadcastTx, err error) {
	return h.services.BroadcastTxSync(packet)
}

func (h *RPCServerContext) broadcastTxCommit(packet []byte) (res *ctypes.ResultBroadcastTxCommit, err error) {
	return h.services.BroadcastTxCommit(packet)
}
