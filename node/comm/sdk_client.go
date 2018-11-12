package comm

import (
	"google.golang.org/grpc"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk/pb"
)

func NewSDKClient() pb.SDKClient {
	address := global.Current.SDKAddress

	// TODO: Why is this insecure????
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("SDK Cannot create client", "err", err)
	}

	client := pb.NewSDKClient(connection)

	return client
}
