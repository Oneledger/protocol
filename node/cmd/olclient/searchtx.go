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

var searchTxCmd = &cobra.Command{
	Use:   "searchtx",
	Short: "Search for transactions",
	Run:   requestTxSearch,
}

var searchTxArgs = &pb.TxSearchRequest{}

func init() {
	RootCmd.AddCommand(searchTxCmd)
	searchTxCmd.Flags().StringVarP(&searchTxArgs.Query, "query", "q", "", "Query for the search")
	searchTxCmd.Flags().BoolVar(&searchTxArgs.Proof, "proof", false, "Include proof for the transactions")
	searchTxCmd.Flags().Int32Var(&searchTxArgs.PerPage, "perage", 10, "Results per page")
	searchTxCmd.Flags().Int32Var(&searchTxArgs.Page, "page", 0, "Page of results")
}

func requestTxSearch(cmd *cobra.Command, args []string) {
	client := comm.NewSDKClient()
	ctx := gcontext.Background()
	reply, err := client.TxSearch(ctx, searchTxArgs)
	if err != nil {
		shared.Console.Error(err)
		os.Exit(1)
	}

	var txResults ctypes.ResultTxSearch
	_, err = serial.Deserialize(reply.Results, &txResults, serial.JSON)
	if err != nil {
		shared.Console.Error(err)
		os.Exit(1)
	}
	shared.Console.Info(txResults)
}
