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
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/abci/types"
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
	Nonce    string `json:"nonce"`
}

func (args *EvmArguments) EVMAccount() (client.EVMAccountRequest, error) {
	return client.EVMAccountRequest{
		Address:  args.From,
		BlockTag: "latest",
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

	var nonce int
	if args.Nonce == "" {
		nonce, err = getNonce(args.From)
	} else {
		nonce, err = strconv.Atoi(args.Nonce)
	}
	if err != nil {
		return client.SendTxRequest{}, err
	}

	return client.SendTxRequest{
		From:     args.From,
		To:       args.To,
		Gas:      int64(args.Gas),
		GasPrice: action.Amount{Currency: "OLT", Value: *gasPriceAmt},
		Amount:   action.Amount{Currency: "OLT", Value: *amt},
		Nonce:    uint64(nonce),
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

	evmExecuteCmd = &cobra.Command{
		Use:   "execute",
		Short: "Execute smart contract or it's code on the OneLedger network",
		RunE:  ExecuteContract,
	}
	evmExecuteArgs = &EvmArguments{}

	evmAccountCmd = &cobra.Command{
		Use:   "account",
		Short: "Get account data from storage",
		RunE:  GetAccount,
	}
	evmAccountArgs = &EvmArguments{}

	evmEstimateCmd = &cobra.Command{
		Use:   "estimate_gas",
		Short: "Estimate execution gas",
		RunE:  EstimateGas,
	}
	evmEstimateArgs = &EvmArguments{}

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

	EVMCmd.AddCommand(evmExecuteCmd)
	setEVMArgs(evmExecuteCmd, evmExecuteArgs)

	EVMCmd.AddCommand(evmAccountCmd)
	setEVMArgs(evmAccountCmd, evmAccountArgs)

	EVMCmd.AddCommand(evmEstimateCmd)
	setEVMArgs(evmEstimateCmd, evmEstimateArgs)

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
	command.Flags().StringVar(&evmArgs.Nonce, "nonce", "", "Nonce of an account. If not set, default will be automatically calculated")
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
		logger.Info("\t Transaction index: ", log.TranscationIndex)
		logger.Info("\t Log index: ", log.LogIndex)
		logger.Info("\t Address: ", log.Address)
		logger.Info("\t Transaction: ", log.TransactionHash)
		logger.Info("\t Block height: ", log.BlockHeight)
		logger.Info("\t Block hash: ", log.BlockHash)
		logger.Info("\t Data: ", log.Data)
		logger.Info("\t Topics: ", log.Topics)
		logger.Info("\t Removed: ", log.Removed)
		logger.Info("\t")
	}
	return nil
}

func getNonce(address keys.Address) (int, error) {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()

	req, err := (&EvmArguments{From: address}).EVMAccount()
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return 0, err
	}

	callResponse, err := fullnode.EVMAccount(req)
	if err != nil {
		logger.Error("error in getting evm call response", err)
		return 0, err
	}
	return int(callResponse.Nonce), nil
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

// EstimateGas used to estimate gas of the deployment code
func EstimateGas(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}

	req, err := evmEstimateArgs.EVMRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	callResponse, err := fullnode.EVMEstimateGas(req)
	if err != nil {
		logger.Error("error in getting evm gas stimate response", err)
		return err
	}
	logger.Info("\t Gas used:", callResponse.GasUsed)
	return nil
}

// CallContract used to invoke function to read from the contract storage
// (mostly for public funcs and simultaion)
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

// ExecuteContract used to deploy/create a contract on the OneLedger protocol
func ExecuteContract(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}

	req, err := evmExecuteArgs.EVMRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	//Prompt for password
	if len(evmExecuteArgs.Password) == 0 {
		evmExecuteArgs.Password = PromptForPassword()
	}

	//Create new Wallet and User Address
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return err
	}

	//Verify User Password
	usrAddress := keys.Address(evmExecuteArgs.From)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, evmExecuteArgs.Password)
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

	if !wallet.Open(usrAddress, evmExecuteArgs.Password) {
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
		PollContractTxResult(ctx, result.Hash.String())
	}

	return nil
}

func PollContractTxResult(ctx *Context, hash string) bool {
	fmt.Println("Checking the contract transaction result...")
	for i := 0; i < queryTxTimes; i++ {
		result, b := checkTransactionResult(ctx, hash, true)
		if result != nil && b == true {
			fmt.Println("Contract transaction is committed!")
			returnData := getReturnData(result.TxResult)
			msgStatus, msgError := getTxStatus(result.TxResult)
			contractAddress := getContractAddress(result.TxResult)
			if len(contractAddress) > 0 {
				fmt.Println("Contract address: ", contractAddress)
			}
			if len(returnData) > 0 {
				fmt.Println("Return data: ", returnData)
			}
			if msgError != "" {
				fmt.Println("Error message: ", msgError)
			}
			fmt.Println("TX status: ", msgStatus)
			fmt.Println("Block height: ", result.Height)
			return true
		}
		time.Sleep(time.Duration(queryTxInternal) * 1000 * time.Millisecond)
	}
	fmt.Println("Transaction failed to be committed in time! Please check later with command: [olclient check_commit] with tx hash")
	return false
}

func getReturnData(resp types.ResponseDeliverTx) (data []byte) {
	for i := range resp.Events {
		evt := resp.Events[i]
		for j := range evt.Attributes {
			attr := evt.Attributes[j]
			if string(attr.Key) == "tx.data" {
				data = attr.Value
			}
		}
	}
	return
}

func getTxStatus(resp types.ResponseDeliverTx) (msgStatus uint64, msgError string) {
	for i := range resp.Events {
		evt := resp.Events[i]
		for j := range evt.Attributes {
			attr := evt.Attributes[j]
			if string(attr.Key) == "tx.status" {
				msgStatus = binary.LittleEndian.Uint64(attr.Value)
			} else if string(attr.Key) == "tx.error" {
				msgError = string(attr.Value)
			}
		}
	}
	return
}

func getContractAddress(resp types.ResponseDeliverTx) (contractAddress ethcmn.Address) {
	for i := range resp.Events {
		evt := resp.Events[i]
		for j := range evt.Attributes {
			attr := evt.Attributes[j]
			if string(attr.Key) == "tx.contract" {
				contractAddress = ethcmn.BytesToAddress(attr.Value)
			}
		}
	}
	return
}
