/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/action"
	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"

	"github.com/Oneledger/protocol/config"
)

type purgeValidatorArgument struct {
	Validator []byte `json:"validator"`
	Admin     []byte `json:"admin"`
}

var purgevalidatorCmd = &cobra.Command{
	Use:   "purgevalidator",
	Short: "Apply a dynamic validator",
	RunE:  purgeValidator,
}

var purgeValidatorArgs = &purgeValidatorArgument{}

func setPurgeValidatorArgs() {
	// Transaction Parameters
	purgevalidatorCmd.Flags().BytesHexVar(&purgeValidatorArgs.Admin, "admin", []byte{nil}, "specify the admin address to use")
	purgevalidatorCmd.Flags().BytesHexVar(&purgeValidatorArgs.Validator, "validator", []byte{nil}, "remove the validator")
}

func init() {
	RootCmd.AddCommand(purgevalidatorCmd)
	setPurgeValidatorArgs()
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func purgeValidator(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	ctx.logger.Debug("Have ApplyValidator Request", "applyValidatorArgs", applyValidatorArgs)

	rootPath, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		return err
	}

	cfg := &config.Server{}

	//Prompt for password
	if len(applyValidatorArgs.Password) == 0 {
		applyValidatorArgs.Password = PromptForPassword()
	}

	//Create new Wallet and User Address
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return err
	}
	//Verify User Password
	usrAddress := keys.Address(applyValidatorArgs.Address)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, applyValidatorArgs.Password)
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
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}

	out, err := fullnode.ApplyValidator(applyValidatorArgs.ClientRequest(currencies.Currencies.GetCurrencySet()))
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

	if !wallet.Open(usrAddress, applyValidatorArgs.Password) {
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
