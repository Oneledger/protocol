package main

import (
  "github.com/Oneledger/protocol/node/comm"
  "github.com/Oneledger/protocol/node/log"
  "github.com/spf13/cobra"
)

var execCmd = &cobra.Command {
  Use: "contract",
  Short: "dry run smart contract",
  Run: RunSmartContract,
}

type ExeArgs struct {
  Address string,
  CallString string,
  CallFrom string,
  SourceCode string,
  Value int,
}

var exeargs ExeArgs = ExeArgs {}

func init() {
  
}
