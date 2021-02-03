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
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/otiai10/copy"
	"path/filepath"
)

const (
	keyStorePath = "keystore/"
)

type LoadTestNode struct {
	clCtx          *client.ExtServiceContext
	cfg            config.Server
	superadmin     keys.Address
	superadminName string
	keypath        string
}

func NewLoadTestNode(superAdmin keys.Address, rootDir string, superAdminName string) LoadTestNode {
	node := LoadTestNode{}
	rootPath, err := filepath.Abs(rootDir)
	if err != nil {
		logger.Fatal(err)
	}
	err = node.cfg.ReadFile(cfgPath(rootPath))
	if err != nil {
		logger.Fatal("failed to read configuration", err)
	}
	clientContext, err := client.NewExtServiceContext(node.cfg.Network.RPCAddress, node.cfg.Network.SDKAddress)
	if err != nil {
		logger.Fatal("error starting rpc client", err)
	}
	keyPath := filepath.Clean(loadTestArgs.superAdminKeyPath)
	keyPathDir := filepath.Dir(loadTestArgs.superAdminKeyPath)
	nodeWallet := filepath.Join(keyPathDir, superAdminName)
	err = copy.Copy(keyPath, nodeWallet)
	if err != nil {
		fmt.Println("Unable to create node", err)
		return node
	}
	node.clCtx = &clientContext
	node.superadmin = superAdmin
	node.superadminName = superAdminName
	node.keypath = nodeWallet
	return node
}

func main() {

	Execute()
}

func cfgPath(dir string) string {
	return filepath.Join(dir, config.FileName)
}
