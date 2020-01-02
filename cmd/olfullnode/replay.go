package main

import (
	"github.com/Oneledger/protocol/app"
	olnode "github.com/Oneledger/protocol/app/node"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type ReplayArgs struct {
	Height int64
}

var replayargs ReplayArgs

var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Start up node (server)",
	RunE:  replay,
}

// Setup the command and flags in Cobra
func init() {
	nodeCmd.AddCommand(replayCmd)

	// Get information to connect to a my tendermint node
	replayCmd.Flags().Int64Var(&replayargs.Height, "height", 0, "the height to replay from")
}

// Start a node to run continously
func replay(cmd *cobra.Command, args []string) error {

	ctx := nodeCtx
	err := ctx.init(rootArgs.rootDir)
	if err != nil {
		return errors.Wrap(err, "failed to initialize config")
	}

	appNodeContext, err := olnode.NewNodeContext(ctx.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create app's node context")
	}

	application, err := app.NewApp(ctx.cfg, appNodeContext)
	if err != nil {
		return errors.Wrap(err, "failed to create new app")
	}

	err = application.Context.Replay(replayargs.Height)

	application.Context.Close()
	return err
}
