/*
	Copyright 2017-2018 OneLedger

	A fullnode for the OneLedger chain. Includes cli arguments to initialize, restart, etc.
*/
package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	go func() {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		log.Println("listen to: " + listener.Addr().String())
		log.Println(http.Serve(listener, nil))
	}()

	Execute() // Pass control to Cobra
}
