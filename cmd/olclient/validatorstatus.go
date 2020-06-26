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
	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/client"
)

var validatorStatusCmd = &cobra.Command{
	Use:   "validatorstatus",
	Short: "Print out validator status",
	Run:   GetValidatorStatus,
}

type ValidatorStatus struct {
	address []byte
	verbose bool
}

var vsArgs *ValidatorStatus = &ValidatorStatus{}

func init() {
	DelegationCmd.AddCommand(validatorStatusCmd)
	validatorStatusCmd.Flags().BytesHexVar(&vsArgs.address, "address", []byte{}, "validator address")
	// TODO: Add implementation in future
	validatorStatusCmd.Flags().BoolVar(&vsArgs.verbose, "verbose", false, "verbose info about validator")
}

func GetValidatorStatus(cmd *cobra.Command, args []string) {
	ctx := NewContext()

	if len(vsArgs.address) == 0 {
		logger.Error("missing address")
		return
	}

	fullnode := ctx.clCtx.FullNodeClient()

	request := client.ValidatorStatusRequest{
		Address: vsArgs.address,
	}
	vs, err := fullnode.ValidatorStatus(request)
	if err != nil {
		logger.Fatal("error in getting validator status", err)
	}
	printValidatorStatus(&vs)
	if vsArgs.verbose {
		// TODO: Add implementation in future
	}
}

func printValidatorStatus(vs *client.ValidatorStatusReply) {
	logger.Info("\t Exists:", vs.Exists)
	logger.Info("\t Height:", vs.Height)
	logger.Info("\t Power:", vs.Power)
	logger.Info("\t Staking:", vs.Staking)
	logger.Info("\t Total delegation amount:", vs.TotalDelegationAmount)
	logger.Info("\t Self delegation amount:", vs.SelfDelegationAmount)
	logger.Info("\t Delegation amount:", vs.DelegationAmount)
}
