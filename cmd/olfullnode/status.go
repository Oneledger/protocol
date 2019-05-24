package main

import (
	"fmt"
	"os"
	"path/filepath"
	"net"
	"net/url"
	
	"github.com/spf13/cobra"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"

	"github.com/kyokomi/emoji"
)

var showStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "show node config file content",
	RunE:  Init,
}

type showStatusArgs struct {
	logger       *log.Logger
	cfg	         *config.Server
	showConfig   bool
}

var showStatusCtx = &showStatusArgs{}

func init() {
	RootCmd.AddCommand(showStatusCmd)

	showStatusCmd.Flags().BoolVarP(&showStatusCtx.showConfig, "verbose", "v", false, "show the node config")
}

func Init(cmd *cobra.Command, args []string) error {
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
	urlRpc, err := url.Parse(cfg.Network.RPCAddress)
	urlP2p, err := url.Parse(cfg.Network.P2PAddress)
	urlSdk, err := url.Parse(cfg.Network.SDKAddress)

	if err != nil {
		return errors.Wrapf(err, "failed to parse url")
	}
	
	host, rpcPort, _ := net.SplitHostPort(urlRpc.Host)
	host, p2pPort, _ := net.SplitHostPort(urlP2p.Host)
	host, sdkPort, _ := net.SplitHostPort(urlSdk.Host)
	
	_, errRpc := net.Listen("tcp", host + ":" + rpcPort)
    if errRpc != nil {
		rpcTaken := emoji.Sprint("RPC Port: ", rpcPort, " on ", host, " :check_mark:")
		fmt.Println(rpcTaken)
    } else {
		rpcAvail := emoji.Sprint("RPC Port: ", rpcPort, " on ", host, " :cross_mark:")
		fmt.Println(rpcAvail)
	}

	_, errP2p := net.Listen("tcp", host + ":" + p2pPort)
    if errP2p != nil {
        p2pTaken := emoji.Sprint("P2P Port: ", p2pPort, " on ", host, " :check_mark:")
		fmt.Println(p2pTaken)
    } else {
		p2pAvail := emoji.Sprint("P2P Port: ", p2pPort, " on ", host, " :cross_mark:")
		fmt.Println(p2pAvail)
	}

	_, errSdk := net.Listen("tcp", host + ":" + sdkPort)
    if errSdk != nil {
        sdkTaken := emoji.Sprint("SDK Port: ", sdkPort, " on ", host, " :check_mark:")
		fmt.Println(sdkTaken)
    } else {
		sdkAvail := emoji.Sprint("SDK Port: ", sdkPort, " on ", host, " :cross_mark:")
		fmt.Println(sdkAvail)
	}
	if errRpc != nil && errP2p != nil && errSdk != nil {
		allHealthy := emoji.Sprint(":clinking_beer_mugs::clinking_beer_mugs::clinking_beer_mugs: Looks all good :clinking_beer_mugs::clinking_beer_mugs::clinking_beer_mugs:")
		fmt.Println(allHealthy)
	}

	return nil  
}

