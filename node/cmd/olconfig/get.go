/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chains.
*/
package main

import (
	"fmt"

	"github.com/Oneledger/protocol/node/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a OneLedger configuration parameter",
	Run:   GetParam,
}

// TODO: typing should be way better, see if cobra can help with this...
type GetArguments struct {
	parameterName string
}

var getParams *GetArguments = &GetArguments{}

func init() {
	RootCmd.AddCommand(getCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	getCmd.Flags().StringVarP(&getParams.parameterName, "parameter", "p", "", "account name")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func GetParam(cmd *cobra.Command, args []string) {
	config.ServerConfig()

	// Don't check it, if it is empty, return empty string
	text := viper.Get(getParams.parameterName)
	if text == nil {
		text = ""
	}

	// Print directly to the screen, so this can be used in shell variables
	fmt.Printf("%s", text)
}
