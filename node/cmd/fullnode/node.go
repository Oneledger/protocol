package main

import (
	"github.com/Oneledger/prototype/node/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/common"
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

	// Set it running
	err := service.Start()
	if err != nil {
		common.Exit(err.Error())
	}

	common.TrapSignal(func() {
		logger.Info("Shutting down")
		service.Stop()
	})
}
