/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/prototype/node/app"
	"github.com/spf13/cobra"
)

var swapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Setup or confirm a currency swap",
	Run:   SwapCurrency,
}

// Arguments to the command
type SwapArguments struct {
	user   string
	to     string
	from   string
	amount []string
}

var swapargs = &SwapArguments{}

func init() {
	RootCmd.AddCommand(swapCmd)

	// Operational Parameters
	//sendCmd.Flags().StringVarP(&app.Current.Transport, "transport", "t", "socket", "transport (socket | grpc)")
	//sendCmd.Flags().StringVarP(&app.Current.Address, "address", "a", "tcp://127.0.0.1:46658", "full address")

	// Transaction Parameters
	swapCmd.Flags().StringVarP(&swapargs.user, "user", "u", "undefined", "user")
	swapCmd.Flags().StringVarP(&swapargs.to, "to", "t", "user", "to user")
	swapCmd.Flags().StringVarP(&swapargs.from, "from", "f", "user", "from user")
	swapCmd.Flags().StringSliceVarP(&swapargs.amount, "amount", "a", []string{"100OLT"}, "coins")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func SwapCurrency(cmd *cobra.Command, args []string) {
	app.Log.Debug("Register Account", "swapargs", swapargs)
}
