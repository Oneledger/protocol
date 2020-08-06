package app

import (
	"fmt"
)

//Follow internalTx.go as a sample

func AddBiddingTXtoQueue(i interface{}) {
	// Add a store similar to the transaction store to external stores .
	// Access that store through app.Context.extStores.
	// Add transaction to the queue from there .
	app := i.(*App)
	fmt.Println("Adding to queue", app.name)
}

func PopBiddingTXfromQueue(i interface{}) {
	//Same as above
	//Pop The TX ,call deliverTX on it
	//Use deliverTxSession to commit or ignore the error
	app := i.(*App)
	fmt.Println("Popping from queue", app.name)
}
