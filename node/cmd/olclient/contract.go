package main

import (
  "github.com/Oneledger/protocol/node/cmd/shared"
  "github.com/Oneledger/protocol/node/log"
  "github.com/spf13/cobra"
  "github.com/Oneledger/protocol/node/olvm/interpreter/vm"
  "github.com/Oneledger/protocol/node/action"
  "github.com/Oneledger/protocol/node/olvm/interpreter/committor"
  //"encoding/gob"
)

var execCmd = &cobra.Command {
  Use: "contract",
  Short: "dry run smart contract",
  Run: RunSmartContract,
}



var contractArgs *shared.ContractArguments = &shared.ContractArguments {}

func init() {
  execCmd.Flags().StringVar(&contractArgs.Address, "address", "", "contract address")
  execCmd.Flags().StringVar(&contractArgs.CallString, "callString", "", "function call string")
  execCmd.Flags().StringVar(&contractArgs.CallFrom, "callFrom", "", "call from")
  execCmd.Flags().StringVar(&contractArgs.SourceCode, "sourceCode", "", "source code")
  execCmd.Flags().IntVar(&contractArgs.Value, "value", 0, "OLT balance attached")
  RootCmd.AddCommand(execCmd)
}

func RunSmartContract(cmd *cobra.Command, args []string) {
  log.Info("This is a dry run command for test your smart contract")
  log.Debug("Have Run Contract Request", "args", args)
  sourceCode := make([]byte,0)
  trasaction := action.Transaction.(action.Contract)

  request := &action.OLVMRequest{
    contractArgs.CallFrom,
    contractArgs.Address,
    contractArgs.CallString,
    contractArgs.Value,
    sourceCode,
    trasaction,
    action.OLVMContext{},
  }
  vm.InitializeClient()
  reply, err := vm.AutoRun(request)
	if err != nil {
		log.Error("Error happens","err",err)
	}
	c := committor.Create()
	log.Info("value returned","return value", reply.Ret)
	log.Info("output updated","output", reply.Out)
	s, _ := c.Commit(reply.Ret, reply.Out)
	log.Info("Transaction created","transaction", s)

}
