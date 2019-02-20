package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"
	"github.com/spf13/cobra"
	//"encoding/gob"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "analyze smart contract",
	Run:   Analyze,
}

func init() {
	analyzeCmd.Flags().StringVar(&contractArgs.Address, "address", "", "contract address")
	analyzeCmd.Flags().StringVar(&contractArgs.CallString, "callString", "", "function call string")
	analyzeCmd.Flags().StringVar(&contractArgs.CallFrom, "callFrom", "", "call from")
	analyzeCmd.Flags().StringVar(&contractArgs.SourceCode, "sourceCode", "", "source code")
	analyzeCmd.Flags().IntVar(&contractArgs.Value, "value", 0, "OLT balance attached")
	RootCmd.AddCommand(analyzeCmd)
}

func Analyze(cmd *cobra.Command, args []string) {
	log.Info("This is a dry run command for test your smart contract")
	log.Debug("Have Run Contract Request", "args", args)

	request := &action.OLVMRequest{
		From:       getDefault(contractArgs.CallFrom, "0x0"),
		Address:    getDefault(contractArgs.Address, "samples://helloworld"),
		CallString: contractArgs.CallString,
		Value:      contractArgs.Value,
		Context:    action.OLVMContext{},
	}
	vm.InitializeClient()
	reply, err := vm.Analyze(request)
	if err != nil {
		log.Error("Error", "err", err)
	}
	log.Info("test contract completed")
	log.Info("value returned", "ret", reply.Ret)
	log.Info("output updated", "output", reply.Out)

}
