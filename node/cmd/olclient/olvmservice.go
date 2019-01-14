package main

import (
  //"github.com/Oneledger/protocol/node/cmd/shared"
  "github.com/Oneledger/protocol/node/log"
  "github.com/spf13/cobra"
  "github.com/Oneledger/protocol/node/olvm/interpreter/vm"
)



var olvmserviceCmd = &cobra.Command {
  Use: "olvmservice",
  Short: "run olvm service",
  Run: RunOLVMService,
}

func init() {
  RootCmd.AddCommand(olvmserviceCmd)
}

func RunOLVMService(cmd *cobra.Command, args []string) {
  log.Info("Launch an OLVM service standalone")
  log.Debug("Have Request", "args", args)
  vm.NewOLVMService().StartService()
}
