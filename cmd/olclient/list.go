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
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/app"
	"github.com/spf13/cobra"
	"path/filepath"
)

type ListArguments struct {
	identityName string
	accountName  string
	validators   bool
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List out Node data",
	Run:   ListNode,
}

var list = &ListArguments{}

func init() {
	RootCmd.AddCommand(listCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	listCmd.Flags().StringVar(&list.identityName, "identity", "", "identity name")
	listCmd.Flags().StringVar(&list.accountName, "account", "", "account name")
	listCmd.Flags().BoolVar(&list.validators, "validators", false, "include validators")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListNode(cmd *cobra.Command, args []string) {

	Ctx := NewContext()



	// test code
	rootPath, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		logger.Error(err)
	}

	cfg := &config.Server{}
	err = cfg.ReadFile(cfgPath(rootPath))
	if err != nil {
		logger.Error("failed to read configuration file at at %s", cfgPath(rootPath))

	}
	nodeCtx, err := app.NewNodeContext(cfg)
	if err != nil {
		logger.Error("failed to create new Node context", err)
	}
	app, err := app.NewApp(cfg ,nodeCtx)
	balanve := app.Context.ValidatorCtx()
	fmt.Println(balanve)
	if err != nil {
		logger.Error("failed to create new app", err)
	}
	//

	req := data.NewRequestFromData("listAccounts", []byte{})
	resp := &data.Response{}
	err = Ctx.clCtx.Query("server.ListAccounts", req, resp)
	if err != nil {
		logger.Error("error in getting all accounts", err)
		return
	}

	var accs = make([]string, 0, 10)
	err = serialize.GetSerializer(serialize.CLIENT).Deserialize(resp.Data, &accs)
	if err != nil {
		logger.Error("error deserializng", err)
		return
	}

	logger.Infof("Accounts on node: %s ", Ctx.cfg.Node.NodeName)
	for _, a := range accs {
		fmt.Println(a)
	}
}
