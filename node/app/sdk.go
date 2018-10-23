package app

import (
	"context"
	"errors"

	"github.com/Oneledger/protocol/node/data"
	status "github.com/Oneledger/protocol/node/err"
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
		return nil, gstatus.Errorf(codes.NotFound, "Account %s not found", accountName)
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
