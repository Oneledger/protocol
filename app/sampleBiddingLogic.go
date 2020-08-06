package app

import (
	"fmt"
)

func AddTXtoQueue(app App) {
	// Add a store similar to the transaction store to external stores .
	// Access that store through app.Context.extStores.
	// Add transaction to the queue from there .
	fmt.Println("Adding to queue")
}

func PopTXfromQueue(app App) {
	//Same as above
	//Pop The TX ,call deliverTX on it
	//
	fmt.Println("Popping from queu")
}
