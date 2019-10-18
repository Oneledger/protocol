/*

 */

package identity

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/identity/internal"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
)

type Job interface {
	DoMyJob(ctx *JobsContext)
	IsMyJobDone(ctx *JobsContext) bool

	IsSufficient() bool
	DoFinalize()

	GetType() string
	GetJobID() string
	IsDone() bool
}

const (
	JobTypeAddSignature = "addSignature"
	JobTypeBTCBroadcast = "btcBroadcast"
)

func makeJob(data []byte, typ string) Job {

	ser := serialize.GetSerializer(serialize.PERSISTENT)

	switch typ {
	case JobTypeAddSignature:
		as := JobAddSignature{}
		ser.Deserialize(data, &as)
		return &as
	}
}

type JobsContext struct {
	service *internal.Service

	trackers *bitcoin.TrackerStore

	BTCPrivKey       *btcec.PrivateKey
	Params           *chaincfg.Params
	ValidatorAddress action.Address
}

func NewJobContext(ctx node.Context, logger *log.Logger,
	router action.Router, tmnode *consensus.Node,
	trackers *bitcoin.TrackerStore) *JobsContext {

	svc := internal.NewService(ctx, logger, router, tmnode)

	return &JobsContext{
		svc,
		trackers,
	}
}
