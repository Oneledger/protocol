package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/Oneledger/protocol/node/status"
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

// Given the name of the identity and the chain type, generate new keys and broadcast this new identity
func (s SDKServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterReply, error) {
	name := request.Identity
	chain := parseChainType(request.Chain)

	// First check if this account already exists, return with error if not
	_, stat := s.App.Accounts.FindNameOnChain(name, chain)
	// If this account already existsm don't let this go through
	if stat != status.MISSING_VALUE {
		return nil, gstatus.Errorf(codes.AlreadyExists, "Identity %s already exists", name)
	}

	privKey, pubKey := id.GenerateKeys(secret(name+chain.String()), true)
	ok := RegisterLocally(s.App, name, chain.String(), chain, pubKey, privKey)
	if !ok {
		return nil, gstatus.Errorf(codes.FailedPrecondition, "Local registration failed")
	}

	packet := action.SignAndPack(
		action.CreateRegisterRequest(
			name,
			chain.String(),
			global.Current.Sequence,
			global.Current.NodeName,
			nil,
			pubKey.Address()))
	comm.Broadcast(packet)

	// TODO: Use proper secret for key generation
	return &pb.RegisterReply{
		Ok:         true,
		Identity:   name,
		PublicKey:  pubKey.Bytes(),
		PrivateKey: privKey.Bytes(),
	}, nil
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
	b := s.App.Utxo.Get(data.DatabaseKey(account.AccountKey()))

	var balance *pb.Balance
	if b != nil {
		balance = &pb.Balance{
			Amount:   b.GetAmountByName("OLT").Amount.Int64(),
			Currency: currencyProtobuf(b.GetAmountByName("OLT").Currency),
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

	packet := action.SignAndPack(send)

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
	// TODO: Use functions in shared package after resolving import cycles
	findBalance := func(key id.AccountKey) (*data.Balance, error) {
		balance := app.Utxo.Get(data.DatabaseKey(key))
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
		action.NewSendInput(pKey, pBalance.GetAmountByName("OLT")),
		action.NewSendInput(cpKey, cpBalance.GetAmountByName("OLT")),
	}

	outputs := []action.SendOutput{
		action.NewSendOutput(pKey, pBalance.GetAmountByName("OLT").Minus(sendAmount)),
		action.NewSendOutput(cpKey, cpBalance.GetAmountByName("OLT").Plus(sendAmount)),
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
	return gstatus.Errorf(codes.NotFound, "Account %s not found", a)
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

func parseChainType(c pb.ChainType) data.ChainType {
	switch c {
	case pb.ChainType_ONELEDGER:
		return data.ONELEDGER
	case pb.ChainType_BITCOIN:
		return data.BITCOIN
	case pb.ChainType_ETHEREUM:
		return data.ETHEREUM
	}
	return data.UNKNOWN
}

func secret(s string) []byte {
	// TODO: proper secret for key generation
	return []byte(s + global.Current.NodeName)
}
