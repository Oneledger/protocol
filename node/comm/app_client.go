package comm

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	client "github.com/tendermint/tendermint/abci/client"
)

// TODO: Why?
var _ *client.Client

// Generic Client interface, allows SetOption
func NewAppClient() client.Client {
	//log.Debug("New Client", "address", global.Current.AppAddress, "transport", global.Current.Transport)

	// TODO: Try multiple times before giving up
	client, err := client.NewClient(global.Current.AppAddress, global.Current.Transport, true)
	if err != nil {
		log.Fatal("Can't start client", "err", err)
	}
	log.Debug("Have Client", "client", client)

	return client
}

// Set an option in the ABCi app directly
func SetOption(key string, value string) {
	log.Debug("Setting Option")

	//client := NewAppClient()
	/*
		options := types.RequestSetOption{
			Key:   key,
			Value: value,
		}
	*/

	//response, err := client.SetOptionSync(options)
	log.Debug("Have Set Option")

	/*
		if err != nil {
			log.Error("SetOption Failed", "err", err, "response", response)
		}
	*/
}
