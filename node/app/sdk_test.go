package app

/*
Tests some some basic RPCs on the SDK
*/

import (
	"context"
	"testing"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func setupRPC() (a *Application, client pb.SDKClient, ctx context.Context, tearDown func()) {
	global.Current.Config.Network.SDKAddress = "http://127.0.0.1:6900"

	a = NewApplication()
	a.Initialize()

	conn, _ := grpc.Dial(global.Current.Config.Network.SDKAddress, grpc.WithInsecure())
	client = pb.NewSDKClient(conn)

	ctx = context.Background()

	tearDown = func() {
		a.Close()
		conn.Close()
	}
	return
}

// TestSDK tests some basic client -> server interaction for each of the methods available
func TestSDK(t *testing.T) {
	application, client, ctx, tearDown := setupRPC()
	defer tearDown()

	t.Run("SDK_Status", func(t *testing.T) {
		in := &pb.StatusRequest{}
		out, _ := client.Status(ctx, in)
		log.Dump("Status", out)
		//assert.Equal(t, out.Ok, true)
	})

	t.Run("SDK_CheckAccount", func(t *testing.T) {
		testName := "bart"
		testPrivateKey, testPublicKey := id.GenerateKeys([]byte(testName), true)
		testChainType := data.ONELEDGER
		testAccount := id.NewAccount(testChainType, testName, testPublicKey, testPrivateKey, "")
		application.Accounts.Add(testAccount)

		in := &pb.CheckAccountRequest{Name: testName}
		expectedOut := &pb.CheckAccountReply{
			Name:       testName,
			Chain:      testChainType.String(),
			AccountKey: id.NewAccountKey(testPublicKey),
			Balance: &pb.Balance{
				Amount:   0,
				Currency: pb.Currency_OLT,
			},
		}
		out, err := client.CheckAccount(ctx, in)
		assert.Equal(t, err, nil)
		assert.Equal(t, out, expectedOut, "Account details should be the same")
	})

	//t.Run("SDK_Send", func(t *testing.T) {})
}
