package cli

import (
	"fmt"
	"os"
)

type CLI struct {
}

type CLIRunner interface {
	Run(args []string)
}

const cmd_createblockchain = "createblockchain"
const cmd_createwallet = "createwallet"
const cmd_getbalance = "getbalance"
const cmd_listaddresses = "listaddresses"
const cmd_printchain = "printchain"
const cmd_reindexutxo = "reindexutxo"
const cmd_send = "send"
const cmd_startnode = "startnode"

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Printf(" %s -address ADDRESS #Create a blockchain and send genesis block reward to ADDRESS\n", cmd_createblockchain)
	fmt.Printf(" %s #Generate a new key-pair and saved it into the wallet file\n", cmd_createwallet)
	fmt.Printf(" %s -address ADDRESS #Get the balance of certain address\n", cmd_getbalance)
	fmt.Printf(" %s #Lists all addresses from the wallet file\n", cmd_listaddresses)
	fmt.Printf(" %s #Print all the blocks of the blockchain\n", cmd_printchain)
	fmt.Printf(" %s #Rebuild the UTXO set\n", cmd_reindexutxo)
	fmt.Printf(" %s - from FROM -to TO -amount AMOUNT -mine #Send AMOUNT of coints from FROM to TO, Mine tot he same node, when -mine is set\n", cmd_send)
	fmt.Printf(" %s -miner ADDRESS #Start a node with ID specified in NODE_ID env. var. -miner enables mining\n", cmd_startnode)
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.ValidateArgs()
	nodeID := os.Getenv("NODE_ID")
	requireNotEmpty(nodeID, "NODE_ID")
	runner := cli.getCliRunner(os.Args[1])
	runner.Run(os.Args[2:])
}

func (cli *CLI) getCliRunner(commandName string) CLIRunner {
	return NewCLIDefaultRunner(cli)
}

func requireNotEmpty(toTest string, label string) {
	if toTest == "" {
		fmt.Printf("%s env. var is not set\n", label)
		os.Exit(1)
	}
}
