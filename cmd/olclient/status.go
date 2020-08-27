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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print out validator status",
	Run:   GetStatus,
}

type Status struct {
	address []byte
	verbose bool
}

var sArgs *Status = &Status{}

func init() {
	DelegationCmd.AddCommand(statusCmd)
	statusCmd.Flags().BytesHexVar(&sArgs.address, "address", []byte{}, "staking/delegator address")
	// TODO: Add implementation in future
	statusCmd.Flags().BoolVar(&sArgs.verbose, "verbose", false, "verbose info about validator")
}

func GetStatus(cmd *cobra.Command, args []string) {
	ctx := NewContext()

	if len(sArgs.address) == 0 {
		logger.Error("missing address")
		return
	}

	fullnode := ctx.clCtx.FullNodeClient()

	request := client.DelegationStatusRequest{
		Address: sArgs.address,
	}
	vs, err := fullnode.DelegationStatus(request)
	if err != nil {
		logger.Fatal("error in getting validator status", err)
	}
	printStatus(&vs)
	if sArgs.verbose {
		// TODO: Add implementation in future
	}
}

func printStatus(vs *client.DelegationStatusReply) {
	logger.Info("\t Balance:", vs.Balance)
	logger.Info("\t Effective delegation amount:", vs.EffectiveDelegationAmount)
	logger.Info("\t Withdrawable amount:", vs.WithdrawableAmount)
	if len(vs.MaturedAmounts) == 0 {
		logger.Info("\t Pending matured amount: empty")
	} else {
		logger.Info("\t Pending matured amount:")
		for _, ma := range vs.MaturedAmounts {
			logger.Infof("\t --- At height: %d, pending amount: %s", ma.Height, ma.Amount.String())
		}
	}
}
