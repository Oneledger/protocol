/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path/filepath"
)

var applyvalidatorCmd = &cobra.Command{
	Use:   "applyvalidator",
	Short: "Apply a dynamic validator",
	RunE:  applyValidator,
}

var applyValidatorArgs *client.ApplyValidatorArguments = &client.ApplyValidatorArguments{}

func init() {
	RootCmd.AddCommand(applyvalidatorCmd)

	// Transaction Parameters
	applyvalidatorCmd.Flags().StringVar(&applyValidatorArgs.Amount, "amount", "0.0", "specify an amount")
	applyvalidatorCmd.Flags().BoolVar(&applyValidatorArgs.Purge, "purge", false, "remove the validator")
	applyvalidatorCmd.Flags().BytesHexVar(&applyValidatorArgs.Address, "address", []byte{}, "address for account to pay the stake and fee")
	applyvalidatorCmd.Flags().BytesBase64Var(&applyValidatorArgs.TmPubKey, "pubkey", []byte{}, "validator pubkey")
	applyvalidatorCmd.Flags().StringVar(&applyValidatorArgs.TmPubKeyType, "pubkeytype", "edd25519", "validator pubkey type, default \"ed25519\", if used pubkey flag, this input should match")
	applyvalidatorCmd.Flags().StringVar(&applyValidatorArgs.Name, "name", "", "name for the validator, default to the node name in config.toml")

}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func applyValidator(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	ctx.logger.Debug("Have ApplyValidator Request", "applyValidatorArgs", applyValidatorArgs)

	rootPath, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		return err
	}

	cfg := &config.Server{}

	err = cfg.ReadFile(cfgPath(rootPath))
	if err != nil {
		return errors.Wrapf(err, "failed to read configuration file at at %s", cfgPath(rootPath))
	}

	// Create message
	resp := &data.Response{}
	err = ctx.clCtx.Query("server.ApplyValidator", *applyValidatorArgs, resp)
	if err != nil {
		ctx.logger.Error("error executing ApplyValidator", err)
		return nil
	}

	packet := resp.Data
	if packet == nil {
		ctx.logger.Error("Error in applying ", resp.ErrorMsg)
		return nil
	}

	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}
	BroadcastStatus(ctx, result)

	return nil

}
