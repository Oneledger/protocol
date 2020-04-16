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
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"

	"github.com/Oneledger/protocol/config"
)

type PurgeValidatorArguments struct {
	Validator []byte `json:"validator"`
	Admin     []byte `json:"admin"`
}

func (args *PurgeValidatorArguments) ClientRequest() client.PurgeValidatorRequest {
	return client.PurgeValidatorRequest{
		Admin:     args.Admin,
		Validator: args.Validator,
	}
}

var purgevalidatorCmd = &cobra.Command{
	Use:   "purgevalidator",
	Short: "Apply a dynamic validator",
	RunE:  purgeValidator,
}

var purgeValidatorArgs = &PurgeValidatorArguments{}

func setPurgeValidatorArgs() {
	// Transaction Parameters
	purgevalidatorCmd.Flags().BytesHexVar(&purgeValidatorArgs.Admin, "admin", []byte{}, "specify the admin address to use")
	purgevalidatorCmd.Flags().BytesHexVar(&purgeValidatorArgs.Validator, "validator", []byte{}, "remove the validator")
}

func init() {
	RootCmd.AddCommand(purgevalidatorCmd)
	setPurgeValidatorArgs()
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func purgeValidator(cmd *cobra.Command, args []string) error {
	ctx := NewContext()
	ctx.logger.Debug("Have PurgeValidator Request", "purgeValidatorArgs", purgeValidatorArgs)

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
	fullnode := ctx.clCtx.FullNodeClient()
	out, err := fullnode.PurgeValidator(purgeValidatorArgs.ClientRequest())
	if err != nil {
		ctx.logger.Error("Error in purging ", err.Error())
		return err
	}

	// Check the reply
	rawTx := &action.RawTx{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(out.RawTx, rawTx)
	if err != nil {
		return errors.New("error de-serializing rawTx")
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(rawTx)
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
