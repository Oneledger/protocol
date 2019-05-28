package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var showStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "show node config file content",
	RunE:  initStatus,
}

type showStatusArgs struct {
	logger     *log.Logger
	cfg        *config.Server
	showConfig bool
}

var showStatusCtx = &showStatusArgs{}

func init() {
	RootCmd.AddCommand(showStatusCmd)

	showStatusCmd.Flags().BoolVarP(&showStatusCtx.showConfig, "verbose", "v", false, "show the node config")
}

func initStatus(cmd *cobra.Command, args []string) error {
	ctx := showStatusCtx

	rootPath, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		return err
	}

	cfg := &config.Server{}
	err = cfg.ReadFile(cfgPath(rootPath))
	if err != nil {
		return errors.Wrapf(err, "failed to read configuration file at at %s", cfgPath(rootPath))
	}

	if ctx.showConfig {
		err := ctx.dumpConfigContent(rootPath, cfg)
		if err != nil {
			return errors.Wrap(err, "failed to initialize config")
		}
	}
	ctx.checkNodes(rootPath, cfg)

	return nil
}

func (ctx *showStatusArgs) dumpConfigContent(rootPath string, cfg *config.Server) error {
	ctx.logger = log.NewLoggerWithPrefix(os.Stdout, "")

	ctx.logger.Dump("[Node]", cfg.Node)
	ctx.logger.Dump("[Network]", cfg.Network)

	return nil
}

func (ctx *showStatusArgs) checkNodes(rootPath string, cfg *config.Server) error {

	okRPC := printPortStatus(cfg.Network.RPCAddress, "RPC")
	okP2P := printPortStatus(cfg.Network.P2PAddress, "P2P")
	okSDK := printPortStatus(cfg.Network.SDKAddress, "SDK")

	if okRPC && okP2P && okSDK {
		fmt.Println("\u2713 Looks all good \u2713")
	}

	return nil
}

func printPortStatus(portAddress string, portType string) bool {
	url, err := url.Parse(portAddress)

	if err != nil {
		fmt.Println(portType, "bad address format \u274C")
		return false
	}

	host, port, _ := net.SplitHostPort(url.Host)

	_, errListen := net.Listen("tcp", host+":"+port)
	if errListen != nil {
		fmt.Println(portType, "Port:", port, "on", host, " \u2713")
		return true
	} else {
		fmt.Println(portType, "Port:", port, "on", host, " \u274C")
		return false
	}
}
