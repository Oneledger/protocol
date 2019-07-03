/*
	Copyright 2017-2018 OneLedger

	A fullnode for the OneLedger chain. Includes cli arguments to initialize, restart, etc.
*/
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {

	go func() {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		r := r1.Intn(100)
		port := 6060 + r
		addr := fmt.Sprintf("localhost:%d", port)
		log.Println(http.ListenAndServe(addr, nil))
	}()

	Execute() // Pass control to Cobra
}
