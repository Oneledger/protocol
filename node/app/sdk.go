package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	status "github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

type SDKServer struct {
	App *Application
}

// Ensure SDKServer implements pb.SDKServer
var _ pb.SDKServer = SDKServer{}

func NewSDKServer(app *Application, port int) (*sdk.Server, error) {
	if port == 0 {
		return nil, errors.New("Failed to start SDK gRPC Server, --sdkrpc was not set")
	}

	s, err := sdk.NewServer(port, &SDKServer{app})
	if err != nil {
		log.Fatal(err.Error())
	}
	return s, nil
}

func (s SDKServer) Status(ctx context.Context, request *pb.StatusRequest) (*pb.StatusReply, error) {
	return &pb.StatusReply{Ok: s.App.SDK.IsRunning()}, nil
}

// CheckAccount returns the balance of a given account ID
func (s SDKServer) CheckAccount(ctx context.Context, request *pb.CheckAccountRequest) (*pb.CheckAccountReply, error) {
	accountName := request.Name
	if accountName == "" {
		return nil, gstatus.Errorf(codes.InvalidArgument, "Account name can't be empty")
	}
	account, err := s.App.Accounts.FindName(accountName)
	if err != status.SUCCESS {
		return nil, errAccountNotFound(accountName)
	}

	// Get balance
	b := s.App.Utxo.Find(data.DatabaseKey(account.AccountKey()))

	var balance *pb.Balance
	if b != nil {
		balance = &pb.Balance{
			Amount:   b.Amount.Amount.Int64(),
			Currency: currencyProtobuf(b.Amount.Currency),
		}
	} else {
		balance = &pb.Balance{Amount: 0, Currency: pb.Currency_OLT}
	}

	return &pb.CheckAccountReply{
		Name:       account.Name(),
		Chain:      account.Chain().String(),
		AccountKey: account.AccountKey().Bytes(),
		Balance:    balance,
	}, nil
}

func (s SDKServer) Send(ctx context.Context, request *pb.SendRequest) (*pb.SendReply, error) {
	findAccount := s.App.Accounts.FindName

	currency := currencyString(request.Currency)
	fee := data.NewCoin(request.Fee, currency)
	gas := data.NewCoin(request.Gas, currency)
	sendAmount := data.NewCoin(request.Amount, currency)

	// Get party & counterparty accounts
	partyAccount, err := findAccount(request.Party)
	if err != status.SUCCESS {
		return nil, errAccountNotFound(request.Party)
	}

	counterPartyAccount, err := findAccount(request.CounterParty)
	if err != status.SUCCESS {
		return nil, errAccountNotFound(request.CounterParty)
	}

	send, errr := prepareSend(partyAccount, counterPartyAccount, sendAmount, fee, gas, s.App)
	if errr != nil {
		return &pb.SendReply{Ok: false, Reason: errr.Error()}, nil
	}

	packet := action.SignAndPack(action.SEND, send)

	// TODO: Include ResultBroadcastTxCommit in SendReply
	comm.Broadcast(packet)

	return &pb.SendReply{
		Ok:     true,
		Reason: "",
	}, nil
}

func prepareSend(
	party id.Account,
	counterParty id.Account,
	sendAmount data.Coin,
	fee data.Coin,
	gas data.Coin,
	app *Application,
) (*action.Send, error) {
	findBalance := func(key id.AccountKey) (*data.Balance, error) {
		balance := app.Utxo.Find(data.DatabaseKey(key))
		if balance == nil {
			return nil, fmt.Errorf("Balance not found for key %x", key)
		}
		return balance, nil
	}

	pKey := party.AccountKey()
	pBalance, errr := findBalance(pKey)
	if errr != nil {
		return nil, fmt.Errorf("Account %s has insufficient balance", party.Name())
	}

	cpKey := counterParty.AccountKey()
	cpBalance, errr := findBalance(cpKey)
	if errr != nil {
		return nil, fmt.Errorf("Account %s has insufficient balance", counterParty.Name())
	}

	inputs := []action.SendInput{
		action.NewSendInput(pKey, pBalance.Amount),
		action.NewSendInput(cpKey, cpBalance.Amount),
	}

	outputs := []action.SendOutput{
		action.NewSendOutput(pKey, pBalance.Amount.Minus(sendAmount)),
		action.NewSendOutput(cpKey, cpBalance.Amount.Plus(sendAmount)),
	}

	return &action.Send{
		Base: action.Base{
			Type:    action.SEND,
			ChainId: ChainId,
			// GetSigners not implemneted
			Signers:  nil,
			Sequence: global.Current.Sequence,
		},
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     fee,
		Gas:     gas,
	}, nil
}

func currencyProtobuf(c data.Currency) pb.Currency {
	switch c.Chain {
	case data.ONELEDGER:
		return pb.Currency_OLT
	case data.BITCOIN:
		return pb.Currency_BTC
	case data.ETHEREUM:
		return pb.Currency_ETH
	}
	// Shouldn't ever reach here
	log.Error("Converting nonexistent currency!!!")
	return pb.Currency_OLT
}

// Functions for returning gRPC errors
func errAccountNotFound(a string) error {
	return gstatus.Errorf(codes.NotFound, "Account %s not found")
}

func currencyString(c pb.Currency) string {
	switch c {
	case pb.Currency_OLT:
		return "OLT"
	case pb.Currency_ETH:
		return "ETH"
	case pb.Currency_BTC:
		return "BTC"
	}

	return "OLT"
}
