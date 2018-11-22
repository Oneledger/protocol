/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	gcontext "context"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/spf13/cobra"
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
		handleError(err)
	}

	out := indentJSON(reply.Results)

	shared.Console.Info(out.String())
}
