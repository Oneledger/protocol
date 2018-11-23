package main

import (
	gcontext "context"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/spf13/cobra"
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
	searchTxCmd.Flags().Int32Var(&searchTxArgs.PerPage, "perpage", 10, "Results per page")
	searchTxCmd.Flags().Int32Var(&searchTxArgs.Page, "page", 0, "Page of results")
}

func requestTxSearch(cmd *cobra.Command, args []string) {
	client := comm.NewSDKClient()
	ctx := gcontext.Background()
	reply, err := client.TxSearch(ctx, searchTxArgs)
	if err != nil {
		handleError(err)
	}

	out := indentJSON(reply.Results)
	shared.Console.Info(out.String())
}
