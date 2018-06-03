/* Copyright 2017 - 2018 OneLedger
 */
package shared

import (
	"os"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

type RegisterArguments struct {
	Identity string
}

// Create a request to register a new identity with the chain
func CreateRegisterRequest(args *RegisterArguments) []byte {
	signers := GetSigners()

	reg := &action.Register{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  app.ChainId,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Identity: args.Identity,
	}

	signed := SignTransaction(action.Transaction(reg))
	packet := PackRequest(signed)

	return packet
}

type SendArguments struct {
	Party        string
	CounterParty string
	Currency     string
	Amount       string
	Gas          string
	Fee          string
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateSendRequest(args *SendArguments) []byte {
	signers := GetSigners()

	conv := convert.NewConvert()

	party := id.Address(conv.GetHash(args.Party))
	counterparty := id.Address(conv.GetHash(args.CounterParty))
	_ = party
	_ = counterparty

	amount := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Amount),
	}
	_ = amount

	gas := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Gas),
	}

	if conv.HasErrors() {
		Console.Error(conv.GetErrors())
		os.Exit(-1)
	}

	// Create base transaction
	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  app.ChainId,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Fee: gas,
		Gas: gas,
	}

	signed := SignTransaction(action.Transaction(send))
	packet := PackRequest(signed)

	return packet
}

// Arguments to the command
type SwapArguments struct {
	Party        string
	CounterParty string
	Amount       string
	Fee          string
	Gas          string // TODO: Not sure this is necessary, unless the chain is like Ethereum
	Currency     string
	Exchange     string
	Excurrency   string
	Nonce        int64
}

// Create a swap request
func CreateSwapRequest(args *SwapArguments) []byte {
	log.Debug("swap args", "args", args)

	// TODO: Need better validation and error handling...

	conv := convert.NewConvert()

	party := id.Address(conv.GetHash(args.Party))
	counterparty := id.Address(conv.GetHash(args.CounterParty))

	// TOOD: a clash with the basic data model
	signers := GetSigners()

	fee := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Fee),
	}

	gas := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Gas),
	}

	amount := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Amount),
	}

	exchange := data.Coin{
		Currency: conv.GetCurrency(args.Excurrency),
		Amount:   conv.GetInt64(args.Exchange),
	}

	if conv.HasErrors() {
		Console.Error(conv.GetErrors())
		os.Exit(-1)
	}

	swap := &action.Swap{
		Base: action.Base{
			Type:     action.SWAP,
			ChainId:  app.ChainId,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Party:        party,
		CounterParty: counterparty,
		Fee:          fee,
		Gas:          gas,
		Amount:       amount,
		Exchange:     exchange,
		Nonce:        args.Nonce,
	}

	signed := SignTransaction(action.Transaction(swap))
	packet := PackRequest(signed)

	return packet
}
