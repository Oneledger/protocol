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
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/data"
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
		log.Fatal("Missing CounterParty information")
	}

	// TODO: Can't convert identities to accounts, this way!
	party := GetAccountKey(args.Party)
	counterParty := GetAccountKey(args.CounterParty)

	if args.Currency == "" || args.Amount == "" {
		log.Fatal("Missing an amount")
	}

	amount := conv.GetCoin(args.Amount, args.Currency)

	// Build up the Inputs
	partyBalance := GetBalance(party)
	counterPartyBalance := GetBalance(counterParty)

	inputs := make([]action.SendInput, 0)
	inputs = append(inputs,
		action.NewSendInput(party, partyBalance),
		action.NewSendInput(counterParty, counterPartyBalance))

	// Build up the outputs
	outputs := make([]action.SendOutput, 0)
	outputs = append(outputs,
		action.NewSendOutput(party, partyBalance.Minus(amount)),
		action.NewSendOutput(counterParty, counterPartyBalance.Plus(amount)))

	fee := conv.GetCoin(args.Fee, args.Currency)
	gas := conv.GetCoin(args.Gas, args.Currency)

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
		log.Fatal("Missing Party information", "args", args)
	}

	// TODO: Can't convert identities to accounts, this way!
	party := GetAccountKey(args.Party)
	zero := GetAccountKey("Zero-OneLedger")

	amount := conv.GetCoin(args.Amount, args.Currency)

	// Build up the Inputs
	zeroBalance := GetBalance(zero)
	partyBalance := GetBalance(party)

	inputs := make([]action.SendInput, 0)
	inputs = append(inputs,
		action.NewSendInput(zero, zeroBalance),
		action.NewSendInput(party, partyBalance))

	// Build up the outputs
	outputs := make([]action.SendOutput, 0)
	outputs = append(outputs,
		action.NewSendOutput(zero, zeroBalance.Minus(amount)),
		action.NewSendOutput(party, partyBalance.Plus(amount)))

	gas := conv.GetCoin(args.Gas, args.Currency)
	fee := conv.GetCoin(args.Fee, args.Currency)

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
			Owner:    zero,
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
	Currency     string
	Fee          string
	Gas          string // TODO: Not sure this is necessary, unless the chain is like Ethereum
	Exchange     string
	Excurrency   string
	Nonce        int64
}

// Create a swap request
func CreateSwapRequest(args *SwapArguments) []byte {
	log.Debug("swap args", "args", args)

	conv := convert.NewConvert()

	partyKey := GetAccountKey(args.Party)
	counterPartyKey := GetAccountKey(args.CounterParty)

	signers := GetSigners()

	fee := conv.GetCoin(args.Fee, "OLT")
	gas := conv.GetCoin(args.Gas, "OLT")
	amount := conv.GetCoin(args.Amount, args.Currency)
	exchange := conv.GetCoin(args.Exchange, args.Excurrency)

	if conv.HasErrors() {
		Console.Error(conv.GetErrors())
		os.Exit(-1)
	}
	account := make(map[data.ChainType]string)
	counterAccount := make(map[data.ChainType]string)

	account[conv.GetChain(args.Currency)] = GetSwapAddress(conv.GetCurrency(args.Currency))
	account[conv.GetChain(args.Excurrency)] = GetSwapAddress(conv.GetCurrency(args.Excurrency))

	party := action.Party{Key: partyKey, Accounts: account}
	counterParty := action.Party{Key: counterPartyKey, Accounts: counterAccount}

	swap := &action.Swap{
		Base: action.Base{
			Type:     action.SWAP,
			ChainId:  app.ChainId,
			Signers:  signers,
			Owner:    partyKey,
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
