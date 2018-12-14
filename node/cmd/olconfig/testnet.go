package main

import (
	"fmt"
	"github.com/Oneledger/protocol/node/global"
	"github.com/tendermint/tendermint/types"
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/node/consensus"
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
	genesis          string
}

var testnetArgs = &testnetConfig{}

var testnetCmd = &cobra.Command{
	Use:   "testnet",
	Short: "Initializes files for a testnet",
	RunE:  runTestnet,
}

func init() {
	RootCmd.AddCommand(testnetCmd)

	testnetCmd.Flags().IntVar(&testnetArgs.numValidators, "validators", 4, "Number of validators to initialize testnet with")
	testnetCmd.Flags().IntVar(&testnetArgs.numNonValidators, "nonvalidators", 0, "Number of non-validators to initialize the testnet with")
	testnetCmd.Flags().StringVar(&testnetArgs.outputDir, "dir", "./oltestnet", "Directory to store initialization files for the testnet")
	testnetCmd.Flags().StringVar(&testnetArgs.genesis, "genesis", "", "Gensis file to use to generate new node key file")
	// TODO: Implement populate persistent peers
}

func runTestnet(cmd *cobra.Command, _ []string) error {
	args := testnetArgs

	if args.numValidators+args.numNonValidators > len(nodeNames) {
		return fmt.Errorf("Don't have enough node names, can't specify more than %d nodes", len(nodeNames))
	}

	if args.genesis != "" {
		genesisdoc, err := types.GenesisDocFromFile(filepath.Join(args.outputDir, args.genesis))
		if err != nil {
			return err
		}
		dir := filepath.Join(args.outputDir, global.Current.ConfigName+"-Node", "consensus", "config")
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		err = genesisdoc.SaveAs(filepath.Join(dir, "genesis.json"))
		if err != nil {
			return err
		}
		// Make node key
		_, err = p2p.LoadOrGenNodeKey(filepath.Join(dir, "node_key.json"))
		if err != nil {
			return err
		}
		// Make private validator file
		pvFile := privval.GenFilePV(filepath.Join(dir, "priv_validator.json"))
		pvFile.Save()
		return nil
	}

	validatorList := make([]consensus.GenesisValidator, testnetArgs.numValidators)

	// Create the GenesisValidator list and its key files priv_validator.json and node_key.json
	for i := 0; i < args.numValidators+args.numNonValidators; i++ {
		isValidator := i < args.numValidators
		nodeName := nodeNames[i]
		nodeDir := filepath.Join(args.outputDir, nodeName+"-Node")
		configDir := filepath.Join(nodeDir, "consensus", "config")
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return err
		}

		// Make node key
		_, err = p2p.LoadOrGenNodeKey(filepath.Join(configDir, "node_key.json"))
		if err != nil {
			return err
		}

		// Make private validator file
		pvFile := privval.GenFilePV(filepath.Join(configDir, "priv_validator.json"))
		pvFile.Save()

		if isValidator {
			validator := consensus.GenesisValidator{
				PubKey: pvFile.GetPubKey(),
				Name:   nodeName,
				Power:  1,
			}
			validatorList[i] = validator
		}
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
	return nil
}
