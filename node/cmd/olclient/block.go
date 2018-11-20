/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	gcontext "context"
	"os"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Look up a block at a given height",
	Run:   requestBlock,
}

var blockArgs = &pb.BlockRequest{}

func init() {
	RootCmd.AddCommand(blockCmd)

	blockCmd.Flags().Int64Var(&blockArgs.Height, "height", 0, "Get the user's height (0 goes to the latest block")
}

func requestBlock(cmd *cobra.Command, args []string) {
	client := comm.NewSDKClient()
	request := &pb.BlockRequest{Height: blockArgs.Height}
	ctx := gcontext.Background()

	reply, err := client.Block(ctx, request)
	if err != nil {
		shared.Console.Error(err)
		os.Exit(1)
	}

	var block ctypes.ResultBlock
	_, err = serial.Deserialize(reply.Results, &block, serial.JSON)
	if err != nil {
		shared.Console.Error(err)
		os.Exit(1)
	}
	shared.Console.Info(block)

}
