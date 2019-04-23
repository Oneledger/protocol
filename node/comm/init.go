package comm

import "github.com/Oneledger/protocol/node/serialize"

var clSerializer serialize.Serializer

func init() {
	clSerializer = serialize.GetSerializer(serialize.CLIENT)
}

const (

	// Broadcast a tx and wait for being commit to block
	BroadcastCommit = "commit"

	// Broadcast a tx and wait for CheckTx response
	// Recommend for smart contract
	BroadcastSync = "sync"

	// Broadcast a tax and return immediately
	BroadcastAsync = "async"
)
