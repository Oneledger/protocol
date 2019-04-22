/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chains.
*/
package main

import (
	"encoding/hex"
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
	validators   bool
}

var list *ListArguments = &ListArguments{}

func init() {
	RootCmd.AddCommand(listCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	listCmd.Flags().StringVar(&list.identityName, "identity", "", "identity name")
	listCmd.Flags().StringVar(&list.accountName, "account", "", "account name")
	listCmd.Flags().BoolVar(&list.validators, "validators", false, "include validators")
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

	ctx := comm.NewClientContext()

	accountRequest := FormatAccountRequest()
	identityRequest := FormatIdentityRequest()

	accounts := ctx.Query("/account", accountRequest)
	identities := ctx.Query("/identity", identityRequest)

	if accounts == nil || identities == nil {
		shared.Console.Warning("No Response from Node for:", string(accountRequest), string(identityRequest))
		return
	}

	nodeName := shared.GetNodeName(ctx)
	printAccountQuery(ctx, nodeName, accounts)
	printIdentityQuery(ctx, nodeName, identities)

	if list.validators == true {
		validators := ctx.Query("/validator", []byte(""))
		if validators != nil {
			printValidatorQuery(nodeName, validators)
		} else {
			shared.Console.Info("Failed to get validator list. Please wait for the next block before running this command.")
		}
	}
}

func printAccountQuery(ctx comm.ClientContext, nodeName string, accountQuery interface{}) {

	accountsI := accountQuery.([]interface{})
	accounts := make([]id.Account, len(accountsI))
	for i := range accountsI {
		a := accountsI[i].(id.Account)
		accounts[i] = a
	}

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
			value := shared.GetBalance(ctx, account.AccountKey())
			if value != nil {
				shared.Console.Info(balance, value.String())
			}
		}
		shared.Console.Info()
	}
}

func printIdentityQuery(ctx comm.ClientContext, nodeName string, idQuery interface{}) {
	identitiesI := idQuery.([]interface{})
	identities := make([]*(id.Identity), len(identitiesI))
	for i := range identitiesI {
		identities[i] = identitiesI[i].(*(id.Identity))
	}

	shared.Console.Info("Identities on", nodeName+":\n")

	for _, identity := range identities {
		printAnIdentity(identity)
	}
}

func printAnIdentity(identity *id.Identity) {
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
	validatorsI := validatorQuery.([]interface{})
	validators := make([]*(id.Identity), len(validatorsI))
	for i := range validatorsI {
		validators[i] = validatorsI[i].(*id.Identity)
	}
	shared.Console.Info("Validators on", nodeName+":\n")

	for _, validator := range validators {
		//printAValidator(validator)
		if validator.Name != "" {
			printAnIdentity(validator)
		}
	}
}

func printAValidator(validator id.Validator) {
	// Right-align fieldnames in console
	address := " Address:"
	pubkey := "  PubKey:"

	shared.Console.Info(address, hex.EncodeToString(validator.Address))
	shared.Console.Info(pubkey, validator.PubKey.String())
	shared.Console.Info()
}
