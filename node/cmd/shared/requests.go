/*
	Copyright 2017 - 2018 OneLedger

	Structures and functions for getting command line arguments, and functions
	to convert these into specific requests.
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

// Prepare a transaction to be issued.
func SignAndPack(ttype action.Type, transaction action.Transaction) []byte {
	signed := action.SignTransaction(transaction)
	packet := action.PackRequest(ttype, signed)

	return packet
}

// Registration
type RegisterArguments struct {
	Identity string
}

// Create a request to register a new identity with the chain
func CreateRegisterRequest(args *RegisterArguments) []byte {
	signers := GetSigners()

	accountKey := GetAccountKey(args.Identity)

	reg := &action.Register{
		Base: action.Base{
			Type:     action.REGISTER,
			ChainId:  app.ChainId,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Identity:   args.Identity,
		NodeName:   global.Current.NodeName,
		AccountKey: accountKey,
	}

	return SignAndPack(action.REGISTER, action.Transaction(reg))
}

// TODO: Get this from the database, need an inquiry from a trusted node?
func GetBalance(account id.AccountKey) data.Coin {
	balance := data.Coin{
		Currency: "OLT",
		Amount:   100,
	}
	return balance
}

type BalanceArguments struct {
}

func CreateBalanceRequest(args *BalanceArguments) []byte {
	return []byte(nil)
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

	if args.Party == "" {
		log.Fatal("Missing Party information")
	}
	if args.CounterParty == "" {
		log.Fatal("Missing Party information")
	}

	// TODO: Can't convert identities to accounts, this way!
	party := conv.GetAccountKey(args.Party)
	counterParty := conv.GetAccountKey(args.CounterParty)

	amount := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Amount),
	}

	// Build up the Inputs
	partyBalance := GetBalance(party)
	counterPartyBalance := GetBalance(counterParty)

	inputs := make([]action.SendInput, 2)
	inputs = append(inputs,
		action.NewSendInput(party, partyBalance),
		action.NewSendInput(counterParty, counterPartyBalance))

	// Build up the outputs
	outputs := make([]action.SendOutput, 2)
	outputs = append(outputs,
		action.NewSendOutput(party, partyBalance.Minus(amount)),
		action.NewSendOutput(counterParty, counterPartyBalance.Plus(amount)))

	gas := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Gas),
	}

	fee := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Fee),
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
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     fee,
		Gas:     gas,
	}

	return SignAndPack(action.SEND, action.Transaction(send))
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateMintRequest(args *SendArguments) []byte {
	signers := GetSigners()

	conv := convert.NewConvert()

	if args.Party == "" {
		log.Fatal("Missing Party information")
	}

	// TODO: Can't convert identities to accounts, this way!
	party := conv.GetAccountKey(args.Party)
	zero := conv.GetAccountKey("Zero")

	amount := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Amount),
	}

	// Build up the Inputs
	zeroBalance := GetBalance(zero)
	partyBalance := GetBalance(party)

	inputs := make([]action.SendInput, 2)
	inputs = append(inputs,
		action.NewSendInput(zero, zeroBalance),
		action.NewSendInput(party, partyBalance))

	// Build up the outputs
	outputs := make([]action.SendOutput, 2)
	outputs = append(outputs,
		action.NewSendOutput(zero, zeroBalance.Minus(amount)),
		action.NewSendOutput(party, partyBalance.Plus(amount)))

	gas := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Gas),
	}

	fee := data.Coin{
		Currency: conv.GetCurrency(args.Currency),
		Amount:   conv.GetInt64(args.Fee),
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
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     fee,
		Gas:     gas,
	}

	return SignAndPack(action.SEND, action.Transaction(send))
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

	party := conv.GetAccountKey(args.Party)
	counterParty := conv.GetAccountKey(args.CounterParty)

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
		CounterParty: counterParty,
		Fee:          fee,
		Gas:          gas,
		Amount:       amount,
		Exchange:     exchange,
		Nonce:        args.Nonce,
	}

	return SignAndPack(action.SWAP, action.Transaction(swap))
}
