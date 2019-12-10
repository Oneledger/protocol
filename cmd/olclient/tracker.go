/*

 */

package main

import (
	"github.com/spf13/cobra"
)

var trackerCmd = &cobra.Command{
	Use:   "tracker",
	Short: "Print out tracker info",
	Run:   TrackerNode,
}

type TrackerReq struct {
	name string
}

var trackerArgs *TrackerReq = &TrackerReq{}

func init() {
	RootCmd.AddCommand(trackerCmd)

	// Transaction Parameters
	trackerCmd.Flags().StringVar(&trackerArgs.name, "tracker_name", "tracker_0", "tracker name")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func TrackerNode(cmd *cobra.Command, args []string) {
	Ctx := NewContext()

	fullnode := Ctx.clCtx.FullNodeClient()
	nodeName, err := fullnode.NodeName()
	if err != nil {
		logger.Fatal(err)
	}

	// assuming we have public key
	bal, err := fullnode.GetTracker(trackerArgs.name)
	if err != nil {
		logger.Fatal("error in getting balance", err)
	}
	printTracker(bal.TrackerData, nodeName)
}

func printTracker(tracker string, nodeName string) {
	logger.Infof("\n %s \n Node Name: %s \n", tracker, nodeName)
}
