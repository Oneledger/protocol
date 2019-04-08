package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/consensus"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

// TODO: Put this in a txt file
var nodeNames = []string{
	"David",
	"Alice",
	"Bob",
	"Carol",
	"Emma",
	"Gary",
	"Harry",
	"Imran",
	"Jason",
	"Kelly",
	"Lisa",
	"Max",
	"Nora",
	"Oliver",
	"Pamela",
	"Quark",
	"Rachel",
	"Sam",
	"Thomas",
	"Ursula",
	"Vladimir",
	"Wendy",
	"Xena",
	"Yuri",
	"Zoey",
}

type testnetConfig struct {
	// Number of validators
	numValidators    int
	numNonValidators int
	outputDir        string
	p2pPort          int
	allowSwap        bool
}

var testnetArgs = &testnetConfig{}

var testnetCmd = &cobra.Command{
	Use:   "devnet",
	Short: "Initializes files for a devnet",
	RunE:  runDevnet,
}

func init() {
	initCmd.AddCommand(testnetCmd)

	testnetCmd.Flags().IntVar(&testnetArgs.numValidators, "validators", 4, "Number of validators to initialize devnet with")
	testnetCmd.Flags().IntVar(&testnetArgs.numNonValidators, "nonvalidators", 0, "Number of non-validators to initialize the devnet with")
	testnetCmd.Flags().StringVarP(&testnetArgs.outputDir, "dir", "o", "./", "Directory to store initialization files for the devnet, default current folder")
	testnetCmd.Flags().BoolVar(&testnetArgs.allowSwap, "enable-swaps", false, "Allow swaps")
}

// Need to maintain a list of nodes and be able to:
// (1) Keep track of all of their P2P addresses including their addresses
// (2) Modify their configurations to have each one have its persistent peer set
type node struct {
	cfg *config.Server
	dir string
	key *p2p.NodeKey
}

func (n node) connectionDetails() string {
	var addr string
	if n.cfg.Network.ExternalP2PAddress == "" {
		addr = n.cfg.Network.P2PAddress
	} else {
		addr = n.cfg.Network.ExternalP2PAddress
	}

	u, _ := url.Parse(addr)
	return fmt.Sprintf("%s@%s", n.key.ID(), u.Host)
}

// This function maintains a running list of ports
// TODO: test if port is available before assigning it by opening a net.Listener
func portGenerator(startingPort int) func() int {
	count := startingPort
	return func() int {
		port := count
		count++
		return port
	}
}

func generateAddress(port int, hasProtocol bool) string {
	var prefix string
	ip := "127.0.0.1"
	protocol := "tcp://"
	if hasProtocol {
		prefix = protocol + ip
	} else {
		prefix = ip
	}

	return fmt.Sprintf("%s:%d", prefix, port)
}

func runDevnet(cmd *cobra.Command, _ []string) error {
	args := testnetArgs
	if args.numValidators+args.numNonValidators > len(nodeNames) {
		return fmt.Errorf("Don't have enough node names, can't specify more than %d nodes", len(nodeNames))
	}
	generatePort := portGenerator(26600)

	validatorList := make([]consensus.GenesisValidator, args.numValidators)
	nodeList := make([]node, args.numValidators+args.numNonValidators)
	persistentPeers := make([]string, args.numNonValidators+args.numValidators)

	// Create the GenesisValidator list and its key files priv_validator_key.json and node_key.json
	for i := 0; i < args.numValidators+args.numNonValidators; i++ {
		isValidator := i < args.numValidators
		nodeName := nodeNames[i]
		nodeDir := filepath.Join(args.outputDir, nodeName+"-Node")
		configDir := filepath.Join(nodeDir, "consensus", "config")
		dataDir := filepath.Join(nodeDir, "consensus", "data")
		nodeDataDir := filepath.Join(nodeDir, "nodedata")

		// Generate new configuration file
		cfg := config.DefaultServerConfig()
		cfg.Node.NodeName = nodeName
		cfg.Network.RPCAddress = generateAddress(generatePort(), true)
		cfg.Network.P2PAddress = generateAddress(generatePort(), true)
		cfg.Network.SDKAddress = generateAddress(generatePort(), true)
		cfg.Network.OLVMAddress = generateAddress(generatePort(), true)

		err := os.MkdirAll(configDir, config.DirPerms)
		if err != nil {
			return err
		}

		err = os.MkdirAll(dataDir, config.DirPerms)
		if err != nil {
			return err
		}

		err = os.MkdirAll(nodeDataDir, config.DirPerms)
		if err != nil {
			return err
		}

		// Make node key
		nodeKey, err := p2p.LoadOrGenNodeKey(filepath.Join(configDir, "node_key.json"))
		if err != nil {
			log.Error("error load or genning node key", "err", err)
			return err
		}

		// Make private validator file
		pvFile := privval.GenFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(dataDir, "priv_validator_state.json"))
		pvFile.Save()

		if isValidator {
			validator := consensus.GenesisValidator{
				PubKey: pvFile.GetPubKey(),
				Name:   nodeName,
				Power:  1,
			}
			validatorList[i] = validator
		}

		// Save the nodes to a list so we can iterate again and
		n := node{cfg, nodeDir, nodeKey}
		nodeList[i] = n
		persistentPeers[i] = n.connectionDetails()
	}

	// Create the non validator nodes

	// Create the genesis file
	genesisDoc := consensus.DefaultGenesisDoc()
	genesisDoc.Validators = validatorList

	for i := 0; i < args.numValidators+args.numNonValidators; i++ {
		nodeName := nodeNames[i]
		nodeDir := filepath.Join(args.outputDir, nodeName+"-Node")
		configDir := filepath.Join(nodeDir, "consensus", "config")
		err := genesisDoc.SaveAs(filepath.Join(configDir, "genesis.json"))
		if err != nil {
			return err
		}
	}

	// Save the files to the node's relevant directory
	generateBTCPort := portGenerator(18831)
	generateETHPort := portGenerator(28101)

	swapNodes := nodeNames[1:4]
	isSwapNode := func(name string) bool {
		for _, nodeName := range swapNodes {
			if nodeName == name {
				return true
			}
		}
		return false
	}
	for _, node := range nodeList {
		node.cfg.P2P.PersistentPeers = persistentPeers
		// Modify the btc and eth ports
		if args.allowSwap && isSwapNode(node.cfg.Node.NodeName) {
			node.cfg.Network.BTCAddress = generateAddress(generateBTCPort(), false)
			node.cfg.Network.ETHAddress = generateAddress(generateETHPort(), false)
		}

		err := node.cfg.SaveFile(filepath.Join(node.dir, config.FileName))
		if err != nil {
			return err
		}
	}

	return nil
}
