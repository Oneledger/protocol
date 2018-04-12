package main

import (
	"github.com/spf13/cobra"
	//"github.com/Oneledger/prototype/node/cmd/fullnode"
	"github.com/Oneledger/prototype/node/app"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/abci/types"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start node",
	Run:   StartNode,
}

func init() {
	RootCmd.AddCommand(nodeCmd)
}

func StartNode(cmd *cobra.Command, args []string) {
	logger.Info("Starting up a Node")
	node := app.NewApplicationContext()
	service = server.NewGRPCServer("unix://data.sock", types.NewGRPCApplication(*node))
	service.SetLogger(logger)
}
