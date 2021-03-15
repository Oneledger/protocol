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
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
)

const (
	keyStorePath    = "keystore/"
	queryTxInternal = 4
	queryTxTimes    = 5
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

var EVMCmd = &cobra.Command{
	Use:   "evm",
	Short: "OneLedger evm",
	Long:  "EVM module for OneLedger chain to execute smart contracts",
}

var DelegationCmd = &cobra.Command{
	Use:   "delegation",
	Short: "OneLedger delegation",
	Long:  "Delegation module for OneLedger chain",
}

var EvidencesCmd = &cobra.Command{
	Use:   "byzantine_fault",
	Short: "OneLedger evidences",
	Long:  "Evidence module for OneLedger chain",
}

var RewardsCmd = &cobra.Command{
	Use:   "rewards",
	Short: "OneLedger rewards",
	Long:  "Rewards module for OneLedger chain",
}

func main() {
	RootCmd.AddCommand(DelegationCmd)
	RootCmd.AddCommand(EvidencesCmd)
	RootCmd.AddCommand(RewardsCmd)
	RootCmd.AddCommand(EVMCmd)
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

func BroadcastStatusSync(ctx *Context, result *ctypes.ResultBroadcastTx) bool {
	if result == nil {
		ctx.logger.Error("Invalid Transaction")
		return false

	} else if result.Code != 0 {
		if result.Code == 200 {
			ctx.logger.Info("Returned Successfully(fullnode query)", result)
			ctx.logger.Info("Result Data", "data", string(result.Data))
			return true
		} else {
			ctx.logger.Error("Syntax, CheckTx Failed", result)
			return false
		}

	} else {
		ctx.logger.Infof("Returned Successfully %#v", result)
		ctx.logger.Info("Result Data", "data", string(result.Data))
		return true
	}
}

func checkTransactionResult(ctx *Context, hash string, prove bool) (*ctypes.ResultTx, bool) {
	fullnode := ctx.clCtx.FullNodeClient()
	result, err := fullnode.CheckCommitResult(hash, prove)
	if err != nil {
		return nil, false
	}
	return &result.Result, true
}

func PollTxResult(ctx *Context, hash string) bool {
	fmt.Println("Checking the transaction result...")
	for i := 0; i < queryTxTimes; i++ {
		result, b := checkTransactionResult(ctx, hash, true)
		if result != nil && b == true {
			fmt.Println("Transaction is committed!")
			return true
		}
		time.Sleep(time.Duration(queryTxInternal) * 1000 * time.Millisecond)
	}
	fmt.Println("Transaction failed to be committed in time! Please check later with command: [olclient check_commit] with tx hash")
	return false
}
