/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"os"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Issue install transaction",
	Run:   IssueInstallRequest,
}

var installArgs *shared.InstallArguments = &shared.InstallArguments{}

func init() {
	RootCmd.AddCommand(installCmd)

	// Transaction Parameters
	installCmd.Flags().StringVar(&installArgs.Owner, "owner", "", "script owner")
	installCmd.Flags().StringVar(&installArgs.Name, "name", "0", "script name")
	installCmd.Flags().StringVar(&installArgs.Version, "version", "", "script version")
	installCmd.Flags().StringVarP(&installArgs.File, "file", "f", "", "script")
	installCmd.Flags().StringVar(&installArgs.Currency, "currency", "OLT", "currency")

	installCmd.Flags().Float64Var(&installArgs.Fee, "fee", 0.0, "include a fee")
	installCmd.Flags().Int64Var(&installArgs.Gas, "gas", 0, "include gas")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueInstallRequest(cmd *cobra.Command, args []string) {
	log.Debug("Have Install Request", "installArgs", installArgs)

	script := shared.MustReadFile(installArgs.File)
	if script == nil {
		shared.Console.Info("CreateInstallRequest no script file", installArgs.File)
	}

	// Create message
	packet := shared.CreateInstallRequest(installArgs, script)
	if packet == nil {
		shared.Console.Error("CreateInstallRequest bad arguments", installArgs)
		os.Exit(-1)
	}

	result := comm.Broadcast(packet)
	BroadcastStatus(result)
}
