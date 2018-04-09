package cli

import (
  "os"
)

type DefaultRunner struct {
  cli *CLI
}

func (runner DefaultRunner) run(args []string) {
  runner.cli.PrintUsage()
  os.Exit(1)
}

func NewCLIDefaultRunner(cli *CLI) CLIRunner{
  return DefaultRunner{cli}
}
