/*
	Copyright 2017 - 2018 OneLedger
*/
package action

import (
	"bytes"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	wire "github.com/tendermint/go-wire"
)

func BroadcastTransaction(ttype Type, transaction Transaction) {
	log.Debug("Broadcast a transaction to the chain")

	// Don't let the death of a client stop the node from running
	defer func() {
		log.Debug("Catching A Panic")
		if r := recover(); r != nil {
			log.Error("Ignoring Client Panic", "r", r)
		}
	}()

	//time.Sleep(10 * time.Second)

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
func SignTransaction(transaction Transaction) Transaction {
	return transaction
}

// Pack a request into a transferable format (wire)
func PackRequest(ttype Type, request Transaction) []byte {
	var base int32

	// Stick a 32 bit integer in front, so that we can identify the struct for deserialization
	buff := new(bytes.Buffer)
	n, err := int(0), error(nil)
	base = int32(ttype)
	wire.WriteInt32(base, buff, &n, &err)
	bytes := buff.Bytes()

	packet, _ := comm.Serialize(request)
	packet = append(bytes, packet...)

	//packet := wire.BinaryBytes(request)
	return packet
}
