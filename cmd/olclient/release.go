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
	"path/filepath"

	"github.com/Oneledger/protocol/action"
	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
)

type ReleaseArguments struct {
	Address  []byte `json:"address"`
	Password string `json:"password"`
}

func (args *ReleaseArguments) ClientRequest() client.ReleaseRequest {
	return client.ReleaseRequest{
		Address: args.Address,
	}
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Perform release to a validator",
	RunE:  releaseExec,
}

var releaseArgs = &ReleaseArguments{}

func init() {
	EvidencesCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().BytesHexVar(&releaseArgs.Address, "address", []byte{}, "address for validator")
	releaseCmd.Flags().StringVar(&releaseArgs.Password, "password", "", "password to access secure wallet")
}

func releaseExec(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	ctx.logger.Debug("Have Release Request", "releaseArgs", releaseArgs)

	rootPath, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		return err
	}

	cfg := &config.Server{}

	//Prompt for password
	if len(releaseArgs.Password) == 0 {
		releaseArgs.Password = PromptForPassword()
	}

	//Create new Wallet and User Address
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return err
	}

	//Verify User Password
	usrAddress := keys.Address(releaseArgs.Address)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, releaseArgs.Password)
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

	out, err := fullnode.Release(releaseArgs.ClientRequest())
	if err != nil {
		ctx.logger.Error("Error in applying ", err.Error())
		return err
	}

	rawTx := &action.RawTx{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(out.RawTx, rawTx)
	if err != nil {
		ctx.logger.Error("failed to deserialize RawTx", err)
		return err
	}

	if !wallet.Open(usrAddress, releaseArgs.Password) {
		ctx.logger.Error("failed to open secure wallet")
		return err
	}

	//Sign Raw "Send" Transaction Using Secure Wallet.
	pub, signature, err := wallet.SignWithAddress(out.RawTx, usrAddress)
	if err != nil {
		ctx.logger.Error("error signing transaction", err)
	}

	signedTx := &action.SignedTx{
		RawTx: *rawTx,
		Signatures: []action.Signature{
			{Signer: pub, Signed: signature},
		},
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if err != nil {
		ctx.logger.Error("failed to serialize signedTx", err)
		return err
	}

	result, err := ctx.clCtx.BroadcastTxSync(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxSync", err)
	}

	if BroadcastStatusSync(ctx, result) {
		PollTxResult(ctx, result.Hash.String())
	}

	return nil
}
