/*

 */

package identity

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity/internal"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
)

type Job interface {
	DoMyJob(ctx *JobsContext, data interface{})
	IsMyJobDone(key keys.PrivateKey, ctx *JobsContext) bool

	IsSufficient() bool
	DoFinalize()

	GetType() string
	GetJobID() string
}

const (
	JobTypeAddSignature = "addSignature"
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
