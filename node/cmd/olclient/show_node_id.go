package main

import (
	"fmt"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/global"
	"github.com/spf13/cobra"
)

var showIP bool

var showNodeIDCMD = &cobra.Command{
	Use:   "show_node_id",
	Short: "Show this node's id",
	RunE:  showNodeID,
}

func init() {
	RootCmd.AddCommand(showNodeIDCMD)
	showNodeIDCMD.Flags().BoolVar(&showIP, "show-ip", showIP, "Show this nodes IP")
}

func showNodeID(_ *cobra.Command, _ []string) error {
	result, err := shared.ShowNodeID(global.Current.Config, showIP)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
