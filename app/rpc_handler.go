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
	"net/url"
	"os"
	"runtime/debug"

	"github.com/google/uuid"

	"github.com/Oneledger/protocol/consensus"
	"github.com/tendermint/tendermint/p2p"

	"github.com/Oneledger/protocol/config"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
)

type RPCServerContext struct {
	nodeName    string
	balances    *balance.Store
	accounts    accounts.Wallet
	currencies  *balance.CurrencyList
	cfg         config.Server
	nodeContext NodeContext

	logger *log.Logger
}

func NewClientHandler(nodeName string, balances *balance.Store, accounts accounts.Wallet,
	currencies *balance.CurrencyList, cfg config.Server, nodeContext NodeContext) *RPCServerContext {

	return &RPCServerContext{nodeName, balances,
		accounts, currencies, cfg, nodeContext,
		log.NewLoggerWithPrefix(os.Stdout, "client_Handler")}
}

// NodeName returns the name of a node. This is useful for displaying it at cmdline.
func (h *RPCServerContext) NodeName(req data.Request, resp *data.Response) error {
	defer h.recoverPanic()

	resp.SetData([]byte(h.nodeName))
	return nil
}

func (h *RPCServerContext) NodeAddress(req data.Request, resp *data.Response) error {
	defer h.recoverPanic()

	address := h.nodeContext.Address()
	resp.SetData([]byte(address))

	return nil
}

func (h *RPCServerContext) NodeID(req data.Request, resp *data.Response) error {
	defer h.recoverPanic()

	configuration, err := consensus.ParseConfig(&h.cfg)
	if err != nil {
		return errors.Wrap(err, "error parsing config")
	}

	nodeKey, err := p2p.LoadNodeKey(configuration.CFG.NodeKeyFile())
	if err != nil {
		return errors.Wrap(err, "error loading node key")
	}

	// silenced error because not present means false
	shouldShowIP, _ := req.GetBool("showIP")

	ip := configuration.CFG.P2P.ExternalAddress
	if shouldShowIP {
		u, err := url.Parse(ip)
		if err != nil {
			return errors.Wrap(err, "error in parsing url")
		}
		resp.JSON(fmt.Sprintf("%s@%s", nodeKey.ID(), u.Host))

	} else {
		resp.JSON(string(nodeKey.ID()))
	}
	return nil
}

type Handler struct {
	balances *balance.Store
	accounts accounts.Wallet
	wallet   accounts.WalletStore

	logger *log.Logger
}

// GetBalance gets the balance of an address
// TODO make it more generic to handle account name and identity
func (h *RPCServerContext) Balance(key []byte, resp *data.Response) error {
	defer h.recoverPanic()

	bal, err := h.balances.Get(key, true)
	if err != nil {
		h.logger.Error("error getting balance", err)
		return errors.Wrap(err, "error getting balance")
	}
	err = resp.SetDataObj(bal)
	if err != nil {
		return errors.Wrap(err, "err serializing for client")
	}

	return err
}

/*
	Account Handlers start here
*/

// AddAccount adds an account to accounts store of the node
func (h *RPCServerContext) AddAccount(acc accounts.Account, resp *data.Response) error {
	defer h.recoverPanic()

	h.logger.Infof("adding account : %#v %s", acc, acc.Address())
	err := h.accounts.Add(acc)
	if err != nil {
		return errors.Wrap(err, "error in adding account to walletstore")
	}

	acc1, err := h.accounts.GetAccount(acc.Address())
	resp.SetDataObj(acc1)

	return nil
}

// DeleteAccount deletes an account from the accounts store of node
func (h *RPCServerContext) DeleteAccount(acc accounts.Account, resp *data.Response) error {
	defer h.recoverPanic()

	err := h.accounts.Delete(acc)
	if err != nil {
		return errors.Wrap(err, "error in deleting account from walletstore")
	}

	return nil
}

// ListAccounts returns a list of all accounts in the accounts store of node
func (h *RPCServerContext) ListAccounts(req data.Request, resp *data.Response) error {
	defer h.recoverPanic()

	accs := h.accounts.Accounts()

	result := make([]string, len(accs))
	for i, a := range accs {
		result[i] = a.String()
	}
	resp.SetDataObj(result)

	return nil
}

func (h *RPCServerContext) SendTx(args client.SendArguments, resp *data.Response) error {
	defer h.recoverPanic()

	send := action.Send{
		From:   keys.Address(args.Party),
		To:     keys.Address(args.CounterParty),
		Amount: args.Amount,
	}

	uuid, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := action.BaseTx{
		Data: send,
		Fee:  fee,
		Memo: uuid.String(),
	}

	pubKey, signed, err := h.accounts.SignWithAddress(tx.Bytes(), send.From)
	if err != nil {
		resp.Error(err.Error())
		return err
	}
	tx.Signatures = []action.Signature{{pubKey, signed}}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return errors.Wrap(err, "err while network serialization")
	}

	resp.SetData(packet)
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

/*
func (r *RPCServerContext) CreateSend(args SendArguments, resp *app.Response) error {
	if args.Party == "" {
		logger.Error("Missing Party argument")
		return errors.New("Missing Party arguments")
	}

	if args.CounterParty == "" {
		logger.Error("Missing CounterParty argument")
		return errors.New("Missing Counterparty arguments")
	}

	party := GetAccountByName(args.Party)
	counterParty := GetAccountByName("Zero")

	if party == nil || counterParty == nil {
		logger.Error("System doesn't recognize the parties", "party", args.Party, "counterParty", args.CounterParty)
		return errors.New("system doesn;t recognize parties")
	}

	if args.Currency == "" || args.Amount == 0.0 {
		logger.Error("Missing an amount argument")
		return errors.New("missing amount")
	}

	amount := balance.NewCoinFromFloat(args.Amount, args.Currency)
	fee := balance.NewCoinFromFloat(args.Fee, "OLT")

	return nil

}
*/
