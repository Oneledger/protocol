package broadcast

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/log"
)

type Service struct {
	logger *log.Logger
	ext    client.ExtServiceContext
}

func NewService(ctx client.ExtServiceContext, logger *log.Logger) *Service {
	return &Service{
		ext:    ctx,
		logger: logger,
	}
}

// Name returns the name of this service. The RPC method will be prefixed with this service name plus a . (e.g. "broadcast.")
func Name() string {
	return "broadcast"
}

// TxAsync returns as soon as the finishes. Returns with a hash
func (svc *Service) TxAsync(req client.BroadcastRequest, reply *client.BroadcastTxReply) error {
	// svc.ext.BroadcastTxAsync
	return nil
}

// TxSync returns when the transaction has been placed inside the mempool
func (svc *Service) TxSync(req client.BroadcastRequest, reply *client.BroadcastTxReply) error {
	// svc.ext.BroadcastTxSync
	return nil
}

// TxCommit returns when the transaction has been committed to a block.
func (svc *Service) TxCommit(req client.BroadcastRequest, reply *client.BroadcastTxReply) error {
	// svc.ext.BroadcastTxCommit
	return nil
}
