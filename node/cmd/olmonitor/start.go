/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"syscall"

	// Import namespace

	"github.com/Oneledger/protocol/node/log"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Startup an OLVM Monitor ",
	Run:   StartMonitor,

	// Shut off all argument handling, and just get everything after start as an argv list.
	DisableFlagParsing: true,
}

type MonitorArgs struct {
	path string   // Path to the active binary
	argv []string // Command line arguments
}

// Declare a shared arguments struct
var arguments = &MonitorArgs{
	argv: make([]string, 0),
}

// Setup the command and flags in Cobra
func init() {
	RootCmd.AddCommand(startCmd)

	//startCmd.Flags().StringVar(&arguments.path, "path", "", "path to execute")
	//startCmd.Flags().StringSliceVar(&arguments.argv, "cmd", arguments.argv, "arguments")
}

// Start a node to run continously
func StartMonitor(cmd *cobra.Command, args []string) {

	// Catch any underlying panics, for now just print out the details properly and stop
	defer func() {
		if r := recover(); r != nil {
			log.Error("Monitor Fatal Panic, shutting down", "r", r)
			debug.PrintStack()
			if service != nil {
				service.Stop()
			}
			os.Exit(-1)
		}
	}()

	log.Dump("Starting up with ", args)
	arguments.path = args[0]
	arguments.argv = append(arguments.argv, args[1:]...)

	if arguments.path == "" {
		log.Fatal("Missing Path")
	}
	log.Dump("Parsed as ", arguments.path, arguments.argv)

	CatchSigterm()

	log.Debug("Waiting forever...")
	rerunProcess(arguments.path, arguments.argv)
}

// Continue to restart the child
func rerunProcess(path string, argv []string) {
	count := 10 // TODO: Needs to be bounded, driven by config, reset by time
	for {
		log.Debug("Executing command", "path", path)
		var stdoutBuf, stderrBuf bytes.Buffer

		command := exec.Command(path, argv...)
		stdoutIn, _ := command.StdoutPipe()
		stderrIn, _ := command.StderrPipe()

		var errStdout, errStderr error
		stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
		stderr := io.MultiWriter(os.Stderr, &stderrBuf)

		if err := command.Start(); err != nil {
			log.Fatal("Invalid Process", "path", path, "argv", argv)
		}

		go func() {
			_, errStdout = io.Copy(stdout, stdoutIn)
		}()

		go func() {
			_, errStderr = io.Copy(stderr, stderrIn)
		}()

		count--
		command.Wait() // Catches SIGCHILD, reaps process correctly
		if count < 1 {
			log.Error("Exhausted Rerun attempts", "path", path, "argv", argv)
			break
		}
		log.Debug("Retrying...")
	}

}

// Run a process in the background
func backgroundProcess(app string) (chan error, *os.Process, error) {
	command := exec.Command(app)
	channel := make(chan error, 1)
	if err := command.Start(); err != nil {
		return nil, nil, err
	}
	go func() {
		channel <- command.Wait()
	}()
	return channel, command.Process, nil
}

// For a running command, manually reap the child before wait?
/*
func reapChild(cmd exec.Command) {
	channel := make(chan os.Signal, 1)
	// Wait on SIGCHILDs
	signal.Notify(channel, syscall.SIGCHLD)
	go func() {
		// Wait for a SIGCHILD to show up
		my_signal := <-channel
		// Redo the wait, to reap properly
		cmd.Wait()
		// Shutdown the child
		signal.Stop(channel)
	}()
}
*/

// A polite way of bring down the service on a SIGTERM
func CatchSigterm() {
	// Catch a SIGTERM and stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range sigs {
			log.Info("Shutting down from Signal", "signal", sig)
			if service != nil {
				service.Stop()
			}
			os.Exit(-1)
		}
	}()

}

/*
// Log all of the global settings
func LogSettings() {
	log.Info("Diagnostics", "Debug", global.Current.Debug, "DisablePasswords", global.Current.DisablePasswords)
	log.Info("Ownership", "NodeName", global.Current.NodeName, "NodeAccountName", global.Current.NodeAccountName,
		"NodeIdentity", global.Current.NodeIdentity)
	log.Info("Locations", "RootDir", global.Current.RootDir)
	log.Info("Addresses", "RpcAddress", global.Current.Config.Network.RPCAddress, "AppAddress", global.Current.AppAddress)
}
*/
