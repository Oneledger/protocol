/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package main

import (
	"fmt"

	"github.com/Oneledger/protocol/client"

	"github.com/spf13/cobra"
)

var showIP bool

var showNodeIDCMD = &cobra.Command{
	Use:   "show_node_id",
	Short: "Show this node's id",
	Run:   showNodeID,
}

func init() {
	RootCmd.AddCommand(showNodeIDCMD)
	showNodeIDCMD.Flags().BoolVar(&showIP, "show-ip", showIP, "Show this nodes IP")
}

func showNodeID(_ *cobra.Command, _ []string) {
	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()

	out, err := fullnode.NodeID(client.NodeIDRequest{ShouldShowIP: showIP})
	if err != nil {
		ctx.logger.Error("failed to get node ID", err)
		return
	}
	fmt.Println(out)
}
