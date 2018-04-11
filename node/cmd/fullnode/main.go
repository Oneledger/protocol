package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var start = &cobra.Command{
	Run:   Start,
	Use:   "start",
	Short: "Start",
	Long:  "Start",
}

func main() {
	fmt.Println("Staring OneLedger Node")
}

func Start(cmd *cobra.Command, args []string) {
	fmt.Println("Start")
}
