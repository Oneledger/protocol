/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"path/filepath"

	"github.com/Oneledger/protocol/action"
	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
)

type WithdrawArguments struct {
	Address  []byte `json:"address"`
	Amount   int64  `json:"amount"`
	Password string `json:"password"`
}

func (args *WithdrawArguments) ClientRequest() client.WithdrawRequest {
	return client.WithdrawRequest{
		Address: args.Address,
		Amount:  *balance.NewAmountFromInt(args.Amount),
	}
}

func init() {
	DelegationCmd.AddCommand(withdrawCmd)
	setWithdrawArgs()
}

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Perform a withdraw of unfrozen stake",
	RunE:  delegateWithdraw,
}

var withdrawArgs = &WithdrawArguments{}

func setWithdrawArgs() {
	// Transaction Parameters
	withdrawCmd.Flags().Int64Var(&withdrawArgs.Amount, "amount", 0, "specify an amount")
	withdrawCmd.Flags().BytesHexVar(&withdrawArgs.Address, "address", []byte{}, "address for account to pay the stake")
	withdrawCmd.Flags().StringVar(&withdrawArgs.Password, "password", "", "password to access secure wallet")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func delegateWithdraw(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	ctx.logger.Debug("Have Withdraw Request", "withdrawArgs", withdrawArgs)

	rootPath, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		return err
	}

	cfg := &config.Server{}

	//Prompt for password
	if len(withdrawArgs.Password) == 0 {
		withdrawArgs.Password = PromptForPassword()
	}

	//Create new Wallet and User Address
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return err
	}

	//Verify User Password
	usrAddress := keys.Address(withdrawArgs.Address)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, withdrawArgs.Password)
	if !authenticated {
		ctx.logger.Error("authentication error", err)
		return err
	}

	err = cfg.ReadFile(cfgPath(rootPath))
	if err != nil {
		return errors.Wrapf(err, "failed to read configuration file at at %s", cfgPath(rootPath))
	}

	// Create message
	fullnode := ctx.clCtx.FullNodeClient()

	out, err := fullnode.Withdraw(withdrawArgs.ClientRequest())
	if err != nil {
		ctx.logger.Error("Error in applying ", err.Error())
		return err
	}

	//Sign Transaction with secure wallet, then append signature to signedTx.
	signedTx := &action.SignedTx{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(out.RawTx, signedTx)
	if err != nil {
		return errors.New("error de-serializing signedTx")
	}

	if !wallet.Open(usrAddress, withdrawArgs.Password) {
		ctx.logger.Error("failed to open secure wallet")
		return errors.New("failed to open secure wallet")
	}

	//Sign Raw "Send" Transaction Using Secure Wallet.
	pub, signature, err := wallet.SignWithAddress(signedTx.RawTx.RawBytes(), usrAddress)
	if err != nil {
		ctx.logger.Error("error signing transaction", err)
		return err
	}

	signedTx.Signatures = append([]action.Signature{{Signer: pub, Signed: signature}}, signedTx.Signatures...)

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if packet == nil || err != nil {
		return errors.New("error serializing packet: " + err.Error())
	}

	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}

	BroadcastStatus(ctx, result)

	return nil
}
