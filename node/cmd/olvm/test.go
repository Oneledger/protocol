package main

import (
  "github.com/Oneledger/protocol/node/cmd/shared"
  "github.com/Oneledger/protocol/node/log"
  "github.com/spf13/cobra"
  "github.com/Oneledger/protocol/node/olvm/interpreter/vm"
  "github.com/Oneledger/protocol/node/action"
  //"encoding/gob"
)

var execCmd = &cobra.Command {
  Use: "test",
  Short: "dry run smart contract",
  Run: Test,
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

func getDefault(originValue, defaultValue string) string {
  if originValue == "" {
    return defaultValue
  } else {
    return originValue
  }
}

func Test(cmd *cobra.Command, args []string) {
  log.Info("This is a dry run command for test your smart contract")
  log.Debug("Have Run Contract Request", "args", args)

  request := &action.OLVMRequest{
    From: getDefault(contractArgs.CallFrom, "0x0"),
    Address:  getDefault(contractArgs.Address,"samples://helloworld"),
    CallString: contractArgs.CallString,
    Value:  contractArgs.Value,
    Context:  action.OLVMContext{},
  }
  vm.InitializeClient()
  reply, err := vm.AutoRun(request)
	if err != nil {
		log.Error("Error","err",err)
	}
  log.Info("test contract completed")
	log.Info("value returned","ret", reply.Ret)
	log.Info("output updated","output", reply.Out)

}
