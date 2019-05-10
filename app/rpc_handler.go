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
	"os"
	"runtime/debug"

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

type RPCServerCtx struct {
	nodeName   string
	balances   *balance.Store
	accounts   accounts.Wallet
	currencies map[string]balance.Currency

	logger *log.Logger
}

func NewClientHandler(nodeName string, balances *balance.Store, accounts accounts.Wallet,
	currencies map[string]balance.Currency) *RPCServerCtx {

	return &RPCServerCtx{nodeName, balances, accounts, currencies,
		log.NewLoggerWithPrefix(os.Stdout, "client_Handler")}
}

// NodeName returns the name of a node. This is useful for displaying it at cmdline.
func (h *RPCServerCtx) NodeName(req data.Request, resp *data.Response) error {
	resp.SetData([]byte(h.nodeName))
	return nil
}

// Currencies returns a list of registered currencies
func (h *RPCServerCtx) Currencies(req data.Request, resp *data.Response) error {
	resp.SetDataObj(h.currencies)
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
func (h *RPCServerCtx) Balance(key []byte, resp *data.Response) error {
	defer h.recoverPanic()

	bal, err := h.balances.Get(key, true)
	resp.SetDataObj(bal)

	fmt.Println(err)
	return err
}

/*
	Account Handlers start here
*/

// AddAccount adds an account to accounts store of the node
func (h *RPCServerCtx) AddAccount(acc accounts.Account, resp *data.Response) error {
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
func (h *RPCServerCtx) DeleteAccount(acc accounts.Account, resp *data.Response) error {
	defer h.recoverPanic()

	err := h.accounts.Delete(acc)
	if err != nil {
		return errors.Wrap(err, "error in deleting account from walletstore")
	}

	return nil
}

// ListAccounts returns a list of all accounts in the accounts store of node
func (h *RPCServerCtx) ListAccounts(req data.Request, resp *data.Response) error {
	defer h.recoverPanic()

	accs := h.accounts.Accounts()
	h.logger.Error("accs", accs)
	resp.SetDataObj(accs)

	return nil
}

func (h *RPCServerCtx) sendTx(args client.SendArguments, resp *data.Response) error {
	send := action.Send{
		From:   keys.Address(args.Party),
		To:     keys.Address(args.CounterParty),
		Amount: args.Amount,
	}

	fee := action.Fee{args.Fee, args.Gas}
	tx := action.BaseTx{
		Data: send,
		Fee:  fee,
	}

	pubKey, signed, err := h.accounts.SignWithAddress(tx.Bytes(), send.From)
	if err != nil {
		resp.Error(err.Error())
		return err
	}
	tx.Signatures = []action.Signature{{pubKey, signed}}

	packet, _ := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)

	resp.SetData(packet)
	return nil
}

/*
	Client handler util methods
*/
// recoverPanic common method to recover from any handler panic
func (h *RPCServerCtx) recoverPanic() {
	if r := recover(); r != nil {
		h.logger.Error("recovering a panic")
		debug.PrintStack()
	}
}

/*
func (r *RPCServerCtx) CreateSend(args SendArguments, resp *app.Response) error {
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
