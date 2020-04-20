/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type PurgeValidatorArguments struct {
	Admin     []byte `json:"admin"`
	Validator []byte `json:"validator"`
	Password  string `json:"password"`
}

func (args *PurgeValidatorArguments) ClientRequest() client.PurgeValidatorRequest {
	return client.PurgeValidatorRequest{
		Admin:     args.Admin,
		Validator: args.Validator,
	}
}

var purgevalidatorCmd = &cobra.Command{
	Use:   "purgevalidator",
	Short: "Purge a validator",
	RunE:  purgeValidator,
}

var purgeValidatorArgs = &PurgeValidatorArguments{}

func setPurgeValidatorArgs() {
	// Transaction Parameters
	purgevalidatorCmd.Flags().BytesHexVar(&purgeValidatorArgs.Admin, "admin", []byte{}, "specify the admin address to use")
	purgevalidatorCmd.Flags().BytesHexVar(&purgeValidatorArgs.Validator, "validator", []byte{}, "remove the validator")
	purgevalidatorCmd.Flags().StringVar(&purgeValidatorArgs.Password, "password", "", "password to access secure wallet")
}

func init() {
	RootCmd.AddCommand(purgevalidatorCmd)
	setPurgeValidatorArgs()
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func purgeValidator(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	ctx.logger.Debug("Have PurgeValidator Request", "purgeValidatorArgs", purgeValidatorArgs)

	//Create new Wallet and User Address
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return err
	}

	//Prompt for password
	if len(purgeValidatorArgs.Password) == 0 {
		purgeValidatorArgs.Password = PromptForPassword()
	}

	//Verify User Password
	usrAddress := keys.Address(purgeValidatorArgs.Admin)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, purgeValidatorArgs.Password)
	if !authenticated {
		ctx.logger.Error("authentication error", err)
		return err
	}

	// Create message
	fullnode := ctx.clCtx.FullNodeClient()
	out, err := fullnode.PurgeValidator(purgeValidatorArgs.ClientRequest())
	if err != nil {
		ctx.logger.Error("Error in purging ", err.Error())
		return err
	}

	rawTx := action.RawTx{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(out.RawTx, &rawTx)
	if err != nil {
		return errors.New("error de-serializing RawTx")
	}

	// Open wallet
	if !wallet.Open(usrAddress, purgeValidatorArgs.Password) {
		ctx.logger.Error("failed to open secure wallet")
		return errors.New("failed to open secure wallet")
	}

	// Sign Raw "Send" Transaction Using Secure Wallet.
	pub, signature, err := wallet.SignWithAddress(out.RawTx, usrAddress)
	if err != nil {
		ctx.logger.Error("error signing transaction", err)
	}

	signatures := []action.Signature{{Signer: pub, Signed: signature}}
	signedTx := &action.SignedTx{
		RawTx:      rawTx,
		Signatures: signatures,
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if err != nil {
		return errors.New("error serializing packet: " + err.Error())
	}

	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}

	BroadcastStatus(ctx, result)

	return nil
}
