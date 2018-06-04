/*
	Copyright 2017 - 2018 OneLedger
*/
package action

import (
	"time"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	wire "github.com/tendermint/go-wire"
)

func SubmitTransaction(transaction Transaction) {
	log.Debug("Send this to the chain")

	// Don't let the death of a client stop the node from running
	defer func() {
		log.Debug("Catching A Panic")
		if r := recover(); r != nil {
			log.Error("Ignoring Client Panic", "r", r)
		}
	}()

	// TODO: Maybe Tendermint isn't ready for transactions...
	// TODO: Can I test this somehow?
	time.Sleep(10 * time.Second)

	packet := SignAndPack(transaction)
	result := comm.Broadcast(packet)

	log.Debug("Submitted Successfully", "result", result)
}

func SignAndPack(transaction Transaction) []byte {
	signed := SignTransaction(transaction)
	packet := PackRequest(signed)

	return packet
}

// SignTransaction with the local keys
func SignTransaction(transaction Transaction) Transaction {
	return transaction
}

// Pack a request into a transferable format (wire)
func PackRequest(request Transaction) []byte {
	packet := wire.BinaryBytes(request)
	return packet
}
