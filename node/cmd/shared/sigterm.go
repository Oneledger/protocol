/*
	Copyright 2017-2018 OneLedger
*/
package shared

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Oneledger/protocol/node/log"
)

// A polite way of bring down the service on a SIGTERM
func CatchSigterm(StopProcess func()) {
	// Catch a SIGTERM and stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range sigs {
			log.Info("Shutting down from Signal", "signal", sig)
			StopProcess()
			/*
				if service != nil {
					service.Stop()
				}
			*/
			os.Exit(-1)
		}
	}()

}
