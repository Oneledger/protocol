/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Issue execute transaction",
	Run:   IssueExecuteRequest,
}

var executeArgs *shared.ExecuteArguments = &shared.ExecuteArguments{}

func init() {
	RootCmd.AddCommand(executeCmd)

	// Transaction Parameters
	executeCmd.Flags().StringVar(&executeArgs.Owner, "owner", "", "script owner")
	executeCmd.Flags().StringVar(&executeArgs.Name, "name", "", "script name")
  executeCmd.Flags().StringVar(&executeArgs.Address, "address", "", "deployed script address")
  executeCmd.Flags().StringVar(&executeArgs.CallString, "callString", "", "call string on the script")
	executeCmd.Flags().StringVar(&executeArgs.Version, "version", "", "script version")
	executeCmd.Flags().StringVar(&executeArgs.Currency, "currency", "OLT", "currency")

	executeCmd.Flags().Float64Var(&executeArgs.Fee, "fee", 0.0, "include a fee")
	executeCmd.Flags().Int64Var(&executeArgs.Gas, "gas", 0, "include gas")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueExecuteRequest(cmd *cobra.Command, args []string) {
	log.Debug("Have Execute Request", "executeArgs", executeArgs)

	// Create message
	packet := shared.CreateExecuteRequest(executeArgs)
	if packet == nil {
		shared.Console.Info("CreateExecuteRequest bad arguments", executeArgs)
		return
	}

	result := comm.Broadcast(packet)
	BroadcastStatus(result)
}
