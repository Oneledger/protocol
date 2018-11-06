/*
	Copyright 2017 - 2018 OneLedger
*/
package action

import (
	"time"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

// Execute a transaction after a specific delay.
// TODO: The node delays in a separate goroutine, but this should really be handled by the consensus engine,
// so that the delay is in the mempool.
func DelayedTransaction(ttype Type, transaction Transaction, waitTime time.Duration) {
	go func(ttype Type, transaction Transaction) {
		time.Sleep(waitTime)
		BroadcastTransaction(ttype, transaction)
	}(ttype, transaction)
}

// Send out the transaction as an async broadcast
func BroadcastTransaction(ttype Type, transaction Transaction) {
	log.Debug("Broadcast a transaction to the chain")

	// Don't let the death of a client stop the node from running
	defer func() {
		if r := recover(); r != nil {
			log.Error("Ignoring Client Panic", "r", r)
		}
	}()

	packet := SignAndPack(ttype, transaction)
	result := comm.Broadcast(packet)

	log.Debug("Submitted Successfully", "result", result)
}

func SignAndPack(ttype Type, transaction Transaction) []byte {
	signed := SignTransaction(transaction)
	packet := PackRequest(ttype, signed)

	return packet
}

// SignTransaction with the local keys
func SignTransaction(transaction Transaction) SignedTransaction {
	packet, err := serial.Serialize(transaction, serial.CLIENT)

	signed := SignedTransaction {transaction, nil}

	if err != nil {
		log.Error("Failed to Serialize packet: ", "error", err)
	} else {
		request := Message(packet)

		response := comm.Query("/signTransaction", request)

		if response == nil {
			log.Warn("Query returned no signature", "request", request)
		} else {
			signed.signature = response.([]byte)
		}
	}

	log.Debug("Transaction signature", "signature", signed.signature)

	return signed
}

// Pack a request into a transferable format (wire)
func PackRequest(ttype Type, request SignedTransaction) []byte {
	packet, err := serial.Serialize(request, serial.CLIENT)
	if err != nil {
		log.Error("Failed to Serialize packet: ", err)
	}

	return packet
}
