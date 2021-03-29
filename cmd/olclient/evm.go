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
	"fmt"
	"strconv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/spf13/cobra"
)

type TransactionLogsArguments struct {
	TransactionHash []byte `json:"transactionHash"`
}

func (args *TransactionLogsArguments) EVMTransactionLogs() (client.EVMTransactionLogsRequest, error) {
	return client.EVMTransactionLogsRequest{
		TransactionHash: args.TransactionHash,
	}, nil
}

type EvmArguments struct {
	From     []byte `json:"from"`
	To       []byte `json:"to,omitempty"`
	Gas      uint64 `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Amount   string `json:"amount"`
	Data     []byte `json:"data,omitempty"`
	Password string `json:"password"`
}

func (args *EvmArguments) EVMAccount() (client.EVMAccountRequest, error) {
	return client.EVMAccountRequest{
		Address: args.From,
	}, nil
}

func (args *EvmArguments) EVMRequest(currencies *balance.CurrencySet) (client.SendTxRequest, error) {
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
		RunE:  CallContract,
	}
	evmCallArgs = &EvmArguments{}

	evmCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create smart contract on the OneLedger network",
		RunE:  CreateContract,
	}
	evmCreateArgs = &EvmArguments{}

	evmAccountCmd = &cobra.Command{
		Use:   "account",
		Short: "Get account data from storage",
		RunE:  GetAccount,
	}
	evmAccountArgs = &EvmArguments{}

	evmTransactionLogsCmd = &cobra.Command{
		Use:   "transaction_logs",
		Short: "Get transaction logs data from storage",
		RunE:  GetTransactionLogs,
	}
	evmTransactionLogsArgs = &TransactionLogsArguments{}
)

func init() {
	EVMCmd.AddCommand(evmCallCmd)
	setEVMArgs(evmCallCmd, evmCallArgs)

	EVMCmd.AddCommand(evmCreateCmd)
	setEVMArgs(evmCreateCmd, evmCreateArgs)

	EVMCmd.AddCommand(evmAccountCmd)
	setEVMArgs(evmAccountCmd, evmAccountArgs)

	EVMCmd.AddCommand(evmTransactionLogsCmd)
	setTransactionLogsArgs(evmTransactionLogsCmd, evmTransactionLogsArgs)
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

func setTransactionLogsArgs(command *cobra.Command, evmArgs *TransactionLogsArguments) {
	command.Flags().BytesHexVar(&evmArgs.TransactionHash, "transactionHash", []byte{}, "The hash of the transaction")
}

func GetTransactionLogs(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()

	req, err := evmTransactionLogsArgs.EVMTransactionLogs()
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	callResponse, err := fullnode.EVMTransactionLogs(req)
	if err != nil {
		logger.Error("error in getting evm call response", err)
		return err
	}
	logger.Info("\t Result: ")
	for _, log := range callResponse.Logs {
		logger.Info("\t Address: ", log.Address)
		logger.Info("\t Transaction: ", log.TransactionHash)
		logger.Info("\t Block height: ", log.BlockHeight)
		logger.Info("\t Block hash: ", log.BlockHash)
		logger.Info("\t Data: ", log.Data)
		logger.Info("\t Topics: ", log.Topics)
	}
	return nil
}

func GetAccount(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()

	req, err := evmAccountArgs.EVMAccount()
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	callResponse, err := fullnode.EVMAccount(req)
	if err != nil {
		logger.Error("error in getting evm call response", err)
		return err
	}
	logger.Info("\t Result: ")
	logger.Info("\t Address: ", callResponse.Address)
	logger.Info("\t Code: ", callResponse.CodeHash)
	logger.Info("\t Balance: ", callResponse.Balance)
	logger.Info("\t Nonce: ", callResponse.Nonce)
	return nil
}

// CallContract used to invoke function to read from the contract storage
// (mostly for public funcs and simultaion to determine gas used)
func CallContract(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}

	req, err := evmCallArgs.EVMRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	callResponse, err := fullnode.EVMCall(req)
	if err != nil {
		logger.Error("error in getting evm call response", err)
		return err
	}
	logger.Info("\t Result:", callResponse.Result)
	return nil
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

	req, err := evmCreateArgs.EVMRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	//Prompt for password
	if len(evmCreateArgs.Password) == 0 {
		evmCreateArgs.Password = PromptForPassword()
	}

	//Create new Wallet and User Address
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return err
	}

	//Verify User Password
	usrAddress := keys.Address(evmCreateArgs.From)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, evmCreateArgs.Password)
	if !authenticated {
		ctx.logger.Error("authentication error", err)
		return err
	}

	//Get Raw "Send" Transaction
	reply, err := fullnode.CreateRawSend(req)
	if err != nil {
		ctx.logger.Error("failed to create SendTx", err)
		return err
	}
	rawTx := &action.RawTx{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(reply.RawTx, rawTx)
	if err != nil {
		ctx.logger.Error("failed to deserialize RawTx", err)
		return err
	}

	if !wallet.Open(usrAddress, evmCreateArgs.Password) {
		ctx.logger.Error("failed to open secure wallet")
		return err
	}

	//Sign Raw "Send" Transaction Using Secure Wallet.
	pub, signature, err := wallet.SignWithAddress(reply.RawTx, usrAddress)
	if err != nil {
		ctx.logger.Error("error signing transaction", err)
	}

	signatures := []action.Signature{{pub, signature}}
	fmt.Println(signatures)
	signedTx := &action.SignedTx{
		RawTx:      *rawTx,
		Signatures: signatures,
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if err != nil {
		ctx.logger.Error("failed to serialize signedTx", err)
		return err
	}

	//Broadcast Transaction
	result, err := ctx.clCtx.BroadcastTxSync(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxSync", err)
	}

	if BroadcastStatusSync(ctx, result) {
		PollTxResult(ctx, result.Hash.String())
	}

	return nil
}
