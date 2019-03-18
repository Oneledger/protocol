package shared

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/spf13/cobra"
)

// NodeFlags takes in a given cobra.Cmd and applies the standard set of flags required to run a Node
func NodeFlags(cmd *cobra.Command) {
	// Get information to connect to a my tendermint node
	cmd.Flags().StringVarP(&global.Current.Config.Network.RPCAddress, "address", "a",
		global.Current.Config.Network.RPCAddress, "consensus address")

	cmd.Flags().BoolVarP(&global.Current.Debug, "debug", "d",
		global.Current.Debug, "Set DEBUG mode")

	cmd.Flags().StringVar(&global.Current.Config.Network.BTCAddress, "btcrpc",
		global.Current.Config.Network.BTCAddress, "bitcoin rpc address")

	cmd.Flags().StringVar(&global.Current.Config.Network.ETHAddress, "ethrpc",
		global.Current.Config.Network.ETHAddress, "ethereum rpc address")

	cmd.Flags().StringVar(&global.Current.Config.Network.SDKAddress, "sdkrpc",
		global.Current.Config.Network.SDKAddress, "Address for SDK RPC Server")

	cmd.Flags().StringArrayVar(&global.Current.PersistentPeers, "persistent_peers", []string{}, "List of persistent peers to connect to")

	// These could be moved to node persistent flags
	cmd.Flags().StringVar(&global.Current.Config.Network.P2PAddress, "p2p", "", "Address to use in P2P network")

	cmd.Flags().StringVar(&global.Current.Seeds, "seeds", "", "List of seeds to connect to")

	cmd.Flags().BoolVar(&global.Current.SeedMode, "seed_mode", false, "List of seeds to connect to")
}
