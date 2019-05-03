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

package client

import (
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
	"os"
)

type Handler struct {
	balances         *storage.ChainState
	accounts         data.Store
	wallet           accounts.WalletStore

	logger 			*log.Logger
}

func NewClientHandler(balances *storage.ChainState, accounts data.Store, wallet accounts.WalletStore) *Handler {
	return &Handler{balances, accounts, wallet,
					log.NewLoggerWithPrefix(os.Stdout, "client_Handler")}
}

func (r *Handler) GetBalance(key []byte, bal *balance.Balance) {

	bal  = r.balances.Get(key, true)
}

func (r *Handler) GetAccount(key []byte, acc *accounts.Account) error {

	// TODO get account by name
	d, err := r.accounts.Get(key)
	if err != nil {
		return err
	}

	acc.FromBytes(d)

	return nil
}

func (r *Handler) AddAccount(acc accounts.Account, resp *data.Response) error {

	err := r.wallet.Add(acc)
	if err != nil {
		return errors.Wrap(err, "error in adding account to walletstore")
	}

	return nil
}

func (r *Handler) DeleteAccount(acc accounts.Account, resp *data.Response) error {

	err := r.wallet.Delete(acc)
	if err != nil {
		return errors.Wrap(err, "error in deleting account from walletstore")
	}

	return nil
}

/*
func (r *Handler) CreateSend(args SendArguments, resp *app.Response) error {
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

