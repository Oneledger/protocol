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

	"github.com/Oneledger/protocol/identity"

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
	fullnode := Ctx.clCtx.FullNodeClient()
	out, err := fullnode.ListValidators()
	if err != nil {
		logger.Error("error in getting all validators", err)
		return
	}

	for _, v := range out.Validators {
		printValidator(v, out.VMap)
	}
	fmt.Println("Height", out.Height)
}

func printValidator(v identity.Validator, vm map[string]bool) {
	isActive := vm[v.Address.String()]
	fmt.Println("Active", isActive)
	fmt.Println("Address", v.Address)
	fmt.Println("StakeAddress", v.StakeAddress)
	fmt.Println("Power", v.Power)
	fmt.Println("Name", v.Name)
	fmt.Println("Staking", v.Staking)

	fmt.Println()

}
