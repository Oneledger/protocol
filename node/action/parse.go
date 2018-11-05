/*
	Copyright 2017-2018 OneLedger

	Parse the incoming transactions

	TODO: switch from individual wire calls, to reading/writing directly to structs
*/
package action

import (
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	wire "github.com/tendermint/go-amino"
)

// Pull out the type, so that the message can be deserialized
func UnpackMessage(message Message) (Type, Message) {
	value, _, err := wire.DecodeInt32(message)
	if err != nil {
		log.Error("UnpackMessage failed", "message", message)
	}
	return Type(value), message[4:]
}

// TODO: Need a better way to handle the polymorphism...
// Parse a message into the appropriate transaction
func Parse(message Message) (Transaction, status.Code) {
	command, body := UnpackMessage(message)

	// TODO: Can I do this with deserialize?
	switch command {

	case SEND:
		action := ParseSend(body)

		return action, status.SUCCESS

	case SWAP:
		action := ParseSwap(body)

		return action, status.SUCCESS

	case EXTERNAL_SEND:
		action := ParseExternalSend(body)

		return action, status.SUCCESS

	case EXTERNAL_LOCK:
		action := ParseExternalLock(body)

		return action, status.SUCCESS

	case PREPARE:
		action := ParsePrepare(body)

		return action, status.SUCCESS

	case COMMIT:
		action := ParseCommit(body)

		return action, status.SUCCESS

	case FORGET:
		action := ParseForget(body)

		return action, status.SUCCESS

	case REGISTER:
		action := ParseRegister(body)

		return action, status.SUCCESS

	case VERIFY:
		action := ParseVerify(body)

		return action, status.SUCCESS

	case PUBLISH:
		action := ParsePublish(body)

		return action, status.SUCCESS
	default:
		log.Error("Unknown transaction", "command", command)
	}

	return nil, status.PARSE_ERROR
}

// Parse a send request
func ParseSend(message Message) *Send {
	log.Debug("Have a Send Request")
	register := &Send{
		Base: Base{Type: SEND},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParseSend", "status", err)
		return nil
	}
	return result.(*Send)
}

// Parse a swap request
func ParseSwap(message Message) *Swap {
	log.Debug("Have a Swap Request")
	register := &Swap{
		Base: Base{Type: SWAP},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParseSwap", "status", err)
		return nil
	}
	return result.(*Swap)
}

// Parse a send request
func ParseExternalSend(message Message) *ExternalSend {
	log.Debug("Have a ExternalSend Request")
	register := &ExternalSend{
		Base: Base{Type: EXTERNAL_SEND},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParseExternalSend", "status", err)
		return nil
	}
	return result.(*ExternalSend)
}

// Parse a send request
func ParseExternalLock(message Message) *ExternalLock {
	log.Debug("Have a ExternalLock Request")
	register := &ExternalLock{
		Base: Base{Type: EXTERNAL_LOCK},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParseExternalLock", "status", err)
		return nil
	}
	return result.(*ExternalLock)
}

// Parse a ready request
func ParsePrepare(message Message) *Prepare {
	log.Debug("Have a Prepare Request")
	register := &Prepare{
		Base: Base{Type: PREPARE},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParsePrepare", "status", err)
		return nil
	}
	return result.(*Prepare)
}

// Parse a ready request
func ParseCommit(message Message) *Commit {
	log.Debug("Have a Commit Request")
	register := &Commit{
		Base: Base{Type: COMMIT},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParseCommit", "status", err)
		return nil
	}
	return result.(*Commit)
}

// Forget the transaction
func ParseForget(message Message) *Forget {
	log.Debug("Have a Forget Request")
	register := &Forget{
		Base: Base{Type: FORGET},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParseForget", "status", err)
		return nil
	}
	return result.(*Forget)
}

// Forget the transaction
func ParseRegister(message Message) *Register {
	log.Debug("Have a Register Request")
	register := &Register{
		Base: Base{Type: REGISTER},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParseRegister", "status", err)
		return nil
	}
	return result.(*Register)
}

// Forget the transaction
func ParseVerify(message Message) *Verify {
	log.Debug("Have a Verify Request", "messsage", message)
	register := &Verify{
		Base: Base{Type: VERIFY},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParseVerify", "status", err)
		return nil
	}
	return result.(*Verify)
}

func ParsePublish(message Message) *Publish {
	log.Debug("Have a Publish Request")
	register := &Publish{
		Base: Base{Type: PUBLISH},
	}

	result, err := serial.Deserialize(message, register, serial.JSON)
	if err != nil {
		log.Error("ParsePublish", "status", err)
		return nil
	}
	return result.(*Publish)
}
