package comm

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk"
	"github.com/Oneledger/protocol/node/sdk/pb"
)

func NewSDKClient() pb.SDKClient {
	address := global.Current.Config.Network.SDKAddress

	// TODO: Why is this insecure????
	connection, err := grpc.Dial(":"+sdk.GetPort(address), grpc.WithInsecure())
	if err != nil {
		log.Fatal("SDK Cannot create client", "err", err)
	}

	client := pb.NewSDKClient(connection)

	return client
}

// Register the request
func SDKRequest(base interface{}) interface{} {
	buffer, err := clSerializer.Serialize(base)
	if err != nil {
		log.Fatal("Serialize Failed", "err", err)
	}
	client := NewSDKClient()
	context := context.Background()

	request := &pb.SDKRequest{
		Parameters: buffer,
	}

	response, err := client.Request(context, request)
	if err != nil {
		log.Fatal("Query Failed", "err", err)
	}

	var prototype interface{}
	err = clSerializer.Deserialize(response.Results, &prototype)
	if err != nil {
		log.Dump("Response was", response.Results)
		log.Fatal("Deserialize Failed", "err", err)
	}

	return prototype
}

// Register the request
func Register(base []byte) {

	client := NewSDKClient()
	context := context.Background()

	request := &pb.SDKRequest{}

	_, err := client.Request(context, request)
	if err != nil {
		log.Fatal("SDK Send failure", "err", err)
	}

	/*
		result, err := serial.Deserialize(reply, prototype, serial.CLIENT)
		if err != nil {
			log.Fatal("Failed to Deserialize", "err", err)
		}
	*/
	return
}
