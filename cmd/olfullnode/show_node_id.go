package main

import (
	"fmt"
	"net/url"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
)

var showIdCmd = &cobra.Command{
	Use:   "show_node_id",
	Short: "show node p2p",
	RunE:  ShowID,
}

type showIdArgs struct {
	showIp bool
}

var showIdCtx = &showIdArgs{}

func init() {
	RootCmd.AddCommand(showIdCmd)

	showIdCmd.Flags().BoolVar(&showIdCtx.showIp, "ip", false, "show the node ip in result")
}

func ShowID(cmd *cobra.Command, args []string) error {

	root := rootArgs.rootDir

	result, err := showNodeID(root, showIdCtx.showIp)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func showNodeID(root string, shouldShowIP bool) (result string, err error) {

	cfg := new(config.Server)
	err = cfg.ReadFile(cfgPath(root))
	if err != nil {
		return "", err
	}
	configuration, err := consensus.ParseConfig(cfg)
	if err != nil {
		return result, err
	}
	nodeKey, err := p2p.LoadNodeKey(configuration.CFG.NodeKeyFile())
	if err != nil {
		return result, err
	}

	ip := configuration.CFG.P2P.ExternalAddress
	if shouldShowIP {
		u, err := url.Parse(ip)
		if err != nil {
			return result, err
		}
		return fmt.Sprintf("%s@%s", nodeKey.ID(), u.Host), nil
	} else {
		return string(nodeKey.ID()), nil
	}
}
