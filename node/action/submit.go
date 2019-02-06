/*
	Copyright 2017 - 2018 OneLedger
*/
package action

import (
	"github.com/Oneledger/protocol/node/id"
	"time"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"

	"github.com/Oneledger/protocol/node/serial"
)

// Execute a transaction after a specific delay.
// TODO: The node delays in a separate goroutine, but this should really be handled by the consensus engine,
// so that the delay is in the mempool.
func DelayedTransaction(transaction Transaction, waitTime time.Duration) {
	go func() {
		time.Sleep(waitTime)
		BroadcastTransaction(transaction, false)
	}()
}

// Send out the transaction as an async broadcast
func BroadcastTransaction(transaction Transaction, sync bool) {
	log.Debug("Broadcast a transaction to the chain")

	// Don't let the death of a client stop the node from running
	defer func() {
		if r := recover(); r != nil {
			log.Error("Ignoring Client Panic", "r", r)
		}
	}()

	packet := SignAndPack(transaction)
	// todo : fix the broadcast result handling
	var result interface{}
	if sync {
		result = comm.Broadcast(packet)
	} else {
		result = comm.BroadcastAsync(packet)
	}

	log.Debug("Submitted Successfully", "result", result)
}

func SignAndPack(transaction Transaction) []byte {
	signed := SignTransaction(transaction)
	packet := PackRequest(signed)

	return packet
}

// SignTransaction with the local keys
func SignTransaction(transaction Transaction) SignedTransaction {
	packet, err := serial.Serialize(transaction, serial.CLIENT)

	signed := SignedTransaction{transaction, nil}

	if err != nil {
		log.Error("Failed to Serialize packet: ", "error", err)
	} else {
		request := Message(packet)

		response := comm.Query("/signTransaction", request)

		if response == nil {
			log.Warn("Query returned no signature", "request", request)
		} else {
			signed.Signatures = response.([]TransactionSignature)
		}
	}

	log.Debug("Transaction signature", "signature", signed.Signatures)

	return signed
}

// Pack a request into a transferable format (wire)
func PackRequest(request SignedTransaction) []byte {
	packet, err := serial.Serialize(request, serial.CLIENT)
	if err != nil {
		log.Error("Failed to Serialize packet: ", err)
	}

	return packet
}

// GetSigners will return the public keys of the signers
func GetSigners(owner []byte) []id.PublicKey {
	publicKey := comm.Query("/accountPublicKey", owner)
	if publicKey == nil {
		return nil
	}

	switch publicKey.(type) {
	case []byte:
		return nil
	}

	return []id.PublicKey{publicKey.(id.PublicKey)}
}
