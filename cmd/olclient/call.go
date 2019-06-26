package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Oneledger/protocol/rpc"
	"github.com/spf13/cobra"
)

var callCmd = &cobra.Command{
	Use:   "call [method] [params]",
	Short: "Make calls to JSONRPC server",
	RunE:  call,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Missing positional arguments")
		}
		params := args[1]
		if !isJSON(params) {
			return errors.New("given params is not a valid JSON string")
		}
		return nil
	},
}

var fullnodeConn string

func init() {
	RootCmd.AddCommand(callCmd)
	callCmd.Flags().StringVar(&fullnodeConn, "rpc_addr", "", "http://127.0.0.1:54321")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func call(cmd *cobra.Command, args []string) error {
	Ctx := NewContext()
	method, rawParams := args[0], args[1]

	if fullnodeConn == "" {
		fullnodeConn = Ctx.cfg.Network.SDKAddress
	}

	var params map[string]interface{}
	err := json.Unmarshal([]byte(rawParams), &params)
	if err != nil {
		return errors.New("failed to marshal json params")
	}

	var reply map[string]interface{}
	client, err := rpc.NewClient(fullnodeConn)
	if err != nil {
		return errors.New("failed to create rpc client")
	}
	err = client.Call(method, params, &reply)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	out, err := json.MarshalIndent(reply, "", "  ")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(out))
	return nil
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
