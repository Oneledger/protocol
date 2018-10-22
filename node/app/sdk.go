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
	accountName := request.AccountName
	if accountName == "" {
		return nil, gstatus.Errorf(codes.InvalidArgument, "accountName can't be empty")
	}
	account, err := s.App.Accounts.FindName(accountName)
	if err != status.SUCCESS {
		return nil, gstatus.Errorf(codes.NotFound, "Account %s not found", accountName)
	}

	// Get balance
	b := s.App.Utxo.Find(data.DatabaseKey(account.AccountKey()))

	// TODO: Don't return string
	var balance string
	if b == nil {
		balance = ""
	} else {
		balance = b.AsString()
	}

	return &pb.CheckAccountReply{
		AccountName: account.Name(),
		Chain:       account.Chain().String(),
		AccountKey:  account.AccountKey().Bytes(),
		Balance:     balance,
	}, nil
}
