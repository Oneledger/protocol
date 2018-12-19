/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chains.
*/
package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List out Node data",
	Run:   ListNode,
}

// TODO: typing should be way better, see if cobra can help with this...
type ListArguments struct {
	identityName string
	accountName  string
}

var list *ListArguments = &ListArguments{}

func init() {
	RootCmd.AddCommand(listCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	listCmd.Flags().StringVar(&list.identityName, "identity", "", "identity name")
	listCmd.Flags().StringVar(&list.accountName, "account", "", "account name")
}

// Format the request into a query structure
func FormatAccountRequest() []byte {
	// Blank user implies all users
	return action.Message("Account=" + list.accountName)
}

func FormatIdentityRequest() []byte {
	// Blank user implies all users
	return action.Message("Identity=" + list.identityName)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListNode(cmd *cobra.Command, args []string) {
	//log.Debug("Checking Account", "account", account)
	accountRequest := FormatAccountRequest()
	identityRequest := FormatIdentityRequest()
	accounts := comm.Query("/account", accountRequest)
	identities := comm.Query("/identity", identityRequest)
	//validators := comm.Query("/validator", []byte(""))

	if accounts == nil || identities == nil {
		shared.Console.Warning("No Response from Node for:", string(accountRequest), string(identityRequest))
		return
	}

	nodeName := shared.GetNodeName()
	printAccountQuery(nodeName, accounts)
	printIdentityQuery(nodeName, identities)
	//printValidatorQuery(nodeName, validators)
}

func printAccountQuery(nodeName string, accountQuery interface{}) {

	accounts := accountQuery.([]id.Account)

	if len(accounts) < 1 {
		shared.Console.Info("No Accounts on", nodeName)
		return
	}

	name := "             Name:"
	balance := "          Balance:"
	accountType := "             Type:"
	accountKey := "              Key:"

	first := true

	for _, account := range accounts {
		if first {
			shared.Console.Info("Accounts on", nodeName+":\n")
			first = false
		}

		shared.Console.Info(name, account.Name())
		shared.Console.Info(accountType, account.Chain().String())
		shared.Console.Info(accountKey, account.AccountKey().String())

		if account.Chain() == data.ONELEDGER {
			value := shared.GetBalance(account.AccountKey())
			if value != nil {
				shared.Console.Info(balance, value.String())
			}
		}
		shared.Console.Info()
	}
}

func printIdentityQuery(nodeName string, idQuery interface{}) {
	identities := idQuery.([]id.Identity)

	shared.Console.Info("Identities on", nodeName+":\n")

	for _, identity := range identities {
		printAnIdentity(identity)
	}
}

func printAnIdentity(identity id.Identity) {
	// Right-align fieldnames in console
	name := "             Name:"
	scope := "            Scope:"
	accountKey := "              Key:"
	tendermintAddress := "TendermintAddress:"
	tendermintPubKey := " TendermintPubKey:"

	var scopeOutput string
	if identity.External {
		scopeOutput = "External"
	} else {
		scopeOutput = "Local"
	}

	shared.Console.Info(name, identity.Name)
	shared.Console.Info(scope, scopeOutput)
	shared.Console.Info(accountKey, identity.AccountKey.String())
	shared.Console.Info(tendermintAddress, identity.TendermintAddress)
	shared.Console.Info(tendermintPubKey, identity.TendermintPubKey)
	shared.Console.Info()
}

func printValidatorQuery(nodeName string, validatorQuery interface{}) {
	//validators := validatorQuery.([]id.ValidatorInfo)
	validators := validatorQuery.([]id.Identity)
	shared.Console.Info("Validators on", nodeName+":\n")

	for _, validator := range validators {
		//printAValidator(validator)
		printAnIdentity(validator)
	}
}

func printAValidator(validator id.ValidatorInfo) {
	// Right-align fieldnames in console
	address := " Address:"
	pubkey := "  PubKey:"

	shared.Console.Info(address, validator.Address)
	shared.Console.Info(pubkey, validator.PubKey)
	shared.Console.Info()
}
