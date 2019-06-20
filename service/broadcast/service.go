package broadcast

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
	"github.com/Oneledger/protocol/serialize"
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
func validateAndSignTx(req client.BroadcastRequest) ([]byte, error) {
	var base action.BaseTx
	err := serialize.GetSerializer(serialize.NETWORK).Deserialize(req.RawTx, &base)
	if err != nil {
		err = rpc.InvalidRequestError("invalid rawTx given")
		return nil, err
	}
	sigs := []action.Signature{{Signer: req.PublicKey, Signed: req.Signature}}
	_, err = action.ValidateBasic(base.Data, base.Fee, base.Memo, sigs)
	if err != nil {
		err = rpc.InvalidRequestError(err.Error())
		return nil, err
	}
	base.Signatures = sigs
	return base.Bytes(), nil
}

// TxAsync returns as soon as the finishes. Returns with a hash
func (svc *Service) TxAsync(req client.BroadcastRequest, reply *client.BroadcastTxReply) error {
	rawSignedTx, err := validateAndSignTx(req)
	if err != nil {
		return err
	}
	result, err := svc.ext.BroadcastTxAsync(rawSignedTx)
	if err != nil {
		return rpc.NewError(rpc.CodeInternalError, err.Error())
	}

	*reply = client.BroadcastTxReply{
		TxHash: []byte(result.Hash),
		Result: *result,
	}
	return nil
}

// TxSync returns when the transaction has been placed inside the mempool
func (svc *Service) TxSync(req client.BroadcastRequest, reply *client.BroadcastTxReply) error {
	rawSignedTx, err := validateAndSignTx(req)
	if err != nil {
		return err
	}
	result, err := svc.ext.BroadcastTxSync(rawSignedTx)
	if err != nil {
		return rpc.NewError(rpc.CodeInternalError, err.Error())
	}

	*reply = client.BroadcastTxReply{
		TxHash: []byte(result.Hash),
		Result: *result,
	}
	return nil
}

// TxCommit returns when the transaction has been committed to a block.
func (svc *Service) TxCommit(req client.BroadcastRequest, reply *client.BroadcastTxCommitReply) error {
	rawSignedTx, err := validateAndSignTx(req)
	if err != nil {
		return err
	}
	result, err := svc.ext.BroadcastTxCommit(rawSignedTx)
	if err != nil {
		return rpc.NewError(rpc.CodeInternalError, err.Error())
	}

	*reply = client.BroadcastTxCommitReply{
		TxHash: []byte(result.Hash),
		Result: *result,
	}
	return nil
}
