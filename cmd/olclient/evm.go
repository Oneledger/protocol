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

package main

import (
	"errors"
	"strconv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/spf13/cobra"
)

type EvmArguments struct {
	From     []byte `json:"from"`
	To       []byte `json:"to,omitempty"`
	Gas      uint64 `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Amount   string `json:"amount"`
	Data     []byte `json:"data,omitempty"`
	Password string `json:"password"`
}

func (args *EvmArguments) ClientRequest(currencies *balance.CurrencySet) (client.SendTxRequest, error) {
	c, ok := currencies.GetCurrencyByName("OLT")
	if !ok {
		return client.SendTxRequest{}, errors.New("OLT currency not supported")
	}
	_, err := strconv.ParseFloat(args.Amount, 64)
	if err != nil {
		return client.SendTxRequest{}, err
	}
	amt := c.NewCoinFromString(padZero(args.Amount)).Amount

	olt, _ := currencies.GetCurrencyByName("OLT")

	_, err = strconv.ParseFloat(args.GasPrice, 64)
	if err != nil {
		return client.SendTxRequest{}, err
	}
	gasPriceAmt := olt.NewCoinFromString(padZero(args.GasPrice)).Amount

	return client.SendTxRequest{
		From:     args.From,
		To:       args.To,
		Gas:      int64(args.Gas),
		GasPrice: action.Amount{Currency: "OLT", Value: *gasPriceAmt},
		Amount:   action.Amount{Currency: "OLT", Value: *amt},
		Data:     args.Data,
	}, nil
}

var (
	evmCallCmd = &cobra.Command{
		Use:   "call",
		Short: "Call smart contract function to fetch storage data",
		Run:   CallContract,
	}
	evmCallArgs = &EvmArguments{}

	evmCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create smart contract on the OneLedger network",
		RunE:  CreateContract,
	}
	evmCreateArgs = &EvmArguments{}
)

func init() {
	EVMCmd.AddCommand(evmCallCmd)
	EVMCmd.AddCommand(evmCreateCmd)

	setEVMArgs(evmCallCmd, evmCallArgs)
	setEVMArgs(evmCreateCmd, evmCreateArgs)
}

func setEVMArgs(command *cobra.Command, evmArgs *EvmArguments) {
	// Transaction Parameters
	command.Flags().BytesHexVar(&evmArgs.From, "from", []byte{}, "The address the transaction is send from")
	command.Flags().BytesHexVar(&evmArgs.To, "to", []byte{}, "The address the transaction is directed to")
	command.Flags().Uint64Var(&evmArgs.Gas, "gas", 90000, "Integer of the gas provided for the transaction execution. It will return unused gas")
	command.Flags().StringVar(&evmArgs.GasPrice, "gasPrice", "0", "Integer of the gasPrice used for each paid gas")
	command.Flags().StringVar(&evmArgs.Amount, "amount", "0", "Integer of the amount sent with this transaction")
	command.Flags().BytesHexVar(&evmArgs.Data, "data", []byte{}, "The compiled code of a contract OR the hash of the invoked method signature and encoded parameters")
	command.Flags().StringVar(&evmArgs.Password, "password", "", "Password to access secure wallet.")
}

// CallContract used to invoke function to read from the contract storage
// (mostly for public funcs and simultaion to determine gas used)
func CallContract(cmd *cobra.Command, args []string) {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return
	}

	req, err := evmCallArgs.ClientRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return
	}

	callResponse, err := fullnode.EVMCall(req)
	if err != nil {
		logger.Error("error in getting evm call response", err)
		return
	}
	logger.Info("\t Result:", callResponse.Result)
}

// CreateContract used to deploy/create a contract on the OneLedger protocol
func CreateContract(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}

	req, err := evmCallArgs.ClientRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	reply, err := fullnode.SendTx(req)
	if err != nil {
		ctx.logger.Error("failed to create SendTx", err)
		return err
	}

	result, err := ctx.clCtx.BroadcastTxSync(reply.RawTx)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxSync", err)
	}

	if BroadcastStatusSync(ctx, result) {
		PollTxResult(ctx, result.Hash.String())
	}

	return nil
}
