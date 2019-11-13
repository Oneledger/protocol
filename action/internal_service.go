/*

 */

package action

import (
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Service struct {
	nodeCtx node.Context

	logger *log.Logger
	router Router

	//only support local client for broadcasting internal txs
	tmrpc *tmclient.Local
}

func NewService(ctx node.Context, logger *log.Logger, router Router, tmnode *consensus.Node) *Service {
	return &Service{
		nodeCtx: ctx,
		logger:  logger,
		router:  router,
		tmrpc:   tmclient.NewLocal(tmnode),
	}
}

func (svc Service) allowedTx(tx RawTx) error {
	h := svc.router.Handler(tx.Type)
	if h == nil {
		return errors.New("transaction type not allowed")
	}
	return nil
}

type InternalBroadcastRequest struct {
	RawTx RawTx `json:"rawTx"`
}

type BroadcastReply struct {
	TxHash Address `json:"txHash"`
	// OK indicates whether this broadcast was a request.
	// For TxSync, it indicates success of CheckTx. Does not guarantee inclusion of a block
	// For TxAsync, it always returns true
	// For TxCommit, it indicates the success of both CheckTx and DeliverTx. If the broadcast fails is false.
	OK     bool   `json:"ok"`
	Height *int64 `json:"height,omitempty"`
	Log    string `json:"log"`
}

func (reply *BroadcastReply) FromResultBroadcastTx(result *ctypes.ResultBroadcastTx) {
	reply.TxHash = Address(result.Hash)
	reply.OK = result.Code == 0
	reply.Height = nil
	reply.Log = result.Log
}

func (reply *BroadcastReply) FromResultBroadcastTxCommit(result *ctypes.ResultBroadcastTxCommit) {
	reply.TxHash = Address(result.Hash)
	reply.OK = result.CheckTx.Code == 0 && result.DeliverTx.Code == 0
	reply.Height = &result.Height
	reply.Log = "check: " + result.CheckTx.Log + ", deliver: " + result.DeliverTx.Log
}

func (svc Service) InternalBroadcast(request InternalBroadcastRequest, reply *BroadcastReply) error {

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
	rawSignedTx := SignedTx{
		RawTx: request.RawTx,
		Signatures: []Signature{Signature{
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
