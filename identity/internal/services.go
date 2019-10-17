/*

 */

package internal

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

type Service struct {
	nodeCtx node.Context

	logger *log.Logger
	router action.Router

	//only support local client for broadcasting internal txs
	tmrpc *tmclient.Local
}

func NewService(ctx node.Context, logger *log.Logger, router action.Router, tmnode *consensus.Node) *Service {
	return &Service{
		nodeCtx: ctx,
		logger:  logger,
		router:  router,
		tmrpc:   tmclient.NewLocal(tmnode),
	}
}

func (svc Service) allowedTx(tx action.RawTx) error {
	h := svc.router.Handler(tx.Type)
	if h == nil {
		return errors.New("transaction type not allowed")
	}
	return nil
}

func (svc Service) InternalBroadcast(request client.InternalBroadcastRequest, reply *client.BroadcastReply) error {

	if err := svc.allowedTx(request.RawTx); err != nil {
		return err
	}

	priKey := svc.nodeCtx.PrivVal()

	h, err := priKey.GetHandler()
	if err != nil {
		return errors.Wrap(err, "wrong node private validator key")
	}
	signed, err := h.Sign(request.RawTx.RawBytes())
	if err != nil {
		return errors.Wrap(err, "signing failed")
	}
	rawSignedTx := action.SignedTx{
		RawTx: request.RawTx,
		Signatures: []action.Signature{action.Signature{
			Signer: h.PubKey(),
			Signed: signed,
		}},
	}

	result, err := svc.tmrpc.BroadcastTxSync(rawSignedTx.SignedBytes())
	if err != nil {
		return errors.Wrap(err, "error broadcast sync")
	}

	reply.FromResultBroadcastTx(result)
	return nil

}
