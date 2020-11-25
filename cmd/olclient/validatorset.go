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
	"sort"

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

	activeList := []identity.Validator{}
	noActiveList := []identity.Validator{}

	for _, v := range out.Validators {
		isActive := out.VMap[v.Address.String()]
		if isActive {
			activeList = append(activeList, v)
		} else {
			noActiveList = append(noActiveList, v)
		}
	}

	// Order validators by power descending
	sort.Slice(activeList, func(i, j int) bool {
		return activeList[i].Power > activeList[j].Power
	})
	sort.Slice(noActiveList, func(i, j int) bool {
		return noActiveList[i].Power > noActiveList[j].Power
	})

	// Print validators
	for _, v := range activeList {
		isFrozen := out.FMap[v.Address.String()]
		printValidator(v, true, isFrozen)
	}
	for _, v := range noActiveList {
		isFrozen := out.FMap[v.Address.String()]
		printValidator(v, false, isFrozen)
	}

	fmt.Println("Height", out.Height)
}

func printValidator(v identity.Validator, isActive bool, isFrozen bool) {
	fmt.Println("Active", isActive)
	fmt.Println("Frozen", isFrozen)
	fmt.Println("Address", v.Address)
	fmt.Println("StakeAddress", v.StakeAddress)
	fmt.Println("Power", v.Power)
	fmt.Println("Name", v.Name)
	fmt.Println("Staking", v.Staking)

	fmt.Println()

}
