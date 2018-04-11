package main

import (
	"fmt"
	"github.com/Oneledger/prototype/node/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/common"
)

var start = &cobra.Command{
	Run:   Start,
	Use:   "start",
	Short: "Start",
	Long:  "Start",
}

var service common.Service

func main() {
	fmt.Println("Staring OneLedger Node")
	node := app.NewApplicationContext()
	_ = node

	// TODO: Just use a generic base, not hooked up yet.
	app := types.NewBaseApplication()

	service = server.NewGRPCServer("unix://data.sock", types.NewGRPCApplication(app))
}

func Start(cmd *cobra.Command, args []string) {
	fmt.Println("Start")
}
