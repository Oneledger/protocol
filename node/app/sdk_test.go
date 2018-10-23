package app

/*
Tests some some basic RPCs on the SDK
*/

import (
	"context"
	"fmt"
	"testing"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func setupRPC() (a *Application, client pb.SDKClient, ctx context.Context, tearDown func()) {
	port := 6969
	global.Current.SDKAddress = port

	a = NewApplication()
	a.Initialize()

	conn, _ := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", port), grpc.WithInsecure())
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
		assert.Equal(t, out.Ok, true)
	})

	t.Run("SDK_CheckAccount", func(t *testing.T) {
		testName := "bart"
		testPrivateKey, testPublicKey := id.GenerateKeys([]byte(testName))
		testChainType := data.ONELEDGER
		testAccount := id.NewAccount(testChainType, testName, testPublicKey, testPrivateKey)
		application.Accounts.Add(testAccount)

		hash1 := id.NewAccountKey(testPublicKey)
		hash2 := id.NewAccountKey(testPublicKey)

		assert.Equal(t, hash1, hash2, "Hasher should generate identical output")

		in := &pb.CheckAccountRequest{Name: testName}
		expectedOut := &pb.CheckAccountReply{
			Name:       testName,
			Chain:      testChainType.String(),
			AccountKey: id.NewAccountKey(testPublicKey),
			Balance:    "",
		}
		out, err := client.CheckAccount(ctx, in)
		assert.Equal(t, err, nil)
		assert.Equal(t, out, expectedOut, "Account details should be the same")
	})
}
