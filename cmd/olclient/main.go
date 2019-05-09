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
	"os"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/log"
)

var logger = log.NewLoggerWithPrefix(os.Stdout, "olclient")

type Context struct {
	logger *log.Logger
	clCtx  *client.Context
}

func NewContext() *Context {
	Ctx := &Context{
		logger: log.NewLoggerWithPrefix(os.Stdout, "olclient"),
	}

	clientContext, err := client.NewContext(client.RPC_ADDRESS)
	if err != nil {
		Ctx.logger.Fatal("error starting rpc client", err)
	}

	Ctx.clCtx = &clientContext
	return Ctx
}

func main() {
	Execute()
}
