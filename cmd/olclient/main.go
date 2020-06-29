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
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	"github.com/spf13/cobra"
)

const (
	keyStorePath = "keystore/"
)

var logger = log.NewLoggerWithPrefix(os.Stdout, "olclient")

type Context struct {
	logger *log.Logger
	clCtx  *client.ExtServiceContext
	cfg    config.Server
}

func NewContext() *Context {
	Ctx := &Context{
		logger: log.NewLoggerWithPrefix(os.Stdout, "olclient"),
	}

	rootPath, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		logger.Fatal(err)
	}

	err = Ctx.cfg.ReadFile(cfgPath(rootPath))
	if err != nil {
		logger.Fatal("failed to read configuration", err)
	}

	clientContext, err := client.NewExtServiceContext(Ctx.cfg.Network.RPCAddress, Ctx.cfg.Network.SDKAddress)
	if err != nil {
		Ctx.logger.Fatal("error starting rpc client", err)
	}

	Ctx.clCtx = &clientContext
	return Ctx
}

var DelegationCmd = &cobra.Command{
	Use:   "delegation",
	Short: "OneLedger delegation",
	Long:  "Delegation module for OneLedger chain",
}

func main() {
	RootCmd.AddCommand(DelegationCmd)
	Execute()
}

func cfgPath(dir string) string {
	return filepath.Join(dir, config.FileName)
}

func PromptForPassword() string {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}

	//New line after password prompt.
	fmt.Println("")

	password := string(bytePassword)
	return strings.TrimSpace(password)
}
