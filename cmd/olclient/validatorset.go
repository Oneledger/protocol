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
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/serialize"

	"github.com/spf13/cobra"
)

var validatorsetCmd = &cobra.Command{
	Use:   "validatorset",
	Short: "List out all validators",
	Run:   ListValidator,
}

func init() {
	RootCmd.AddCommand(validatorsetCmd)

}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListValidator(cmd *cobra.Command, args []string) {

	Ctx := NewContext()

	req := data.NewRequestFromData("listValidators", []byte{})
	resp := &data.Response{}
	err := Ctx.clCtx.Query("server.ListValidators", req, resp)
	if err != nil {
		logger.Error("error in getting all validators", err)
		return
	}

	var validators = make([]identity.Validator, 0)

	err = serialize.GetSerializer(serialize.CLIENT).Deserialize(resp.Data, &validators)
	if err != nil {
		logger.Error("error deserializng", err)
		return
	}

	logger.Infof("Validators on node: %s ", Ctx.cfg.Node.NodeName)
	for _, v := range validators {
		printValidator(v)
	}
}

func printValidator (v identity.Validator) {
	fmt.Println("Address", v.Address)
	fmt.Println("StakeAddress", v.StakeAddress)
	fmt.Println("Power", v.Power)
	fmt.Println("Name", v.Name)
	fmt.Println("Staking", v.Staking)
	fmt.Println()

}
