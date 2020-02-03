package ethereum

import (
	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (svc *Service) PrepareOLTERC20Lock(req *OLTERC20LockRequest, out *OLTReply) error {
	erc20lock := eth.ERC20Lock{
		Locker: req.Address,
		ETHTxn: req.RawTx,
	}

	data, err := erc20lock.Marshal()
	if err != nil {
		svc.logger.Error(err, codes.ErrPreparingErc20OLTLock.ErrorMsg())
		return codes.ErrPreparingErc20OLTLock
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{req.Fee, req.Gas}
	tx := &action.RawTx{
		Type: action.ERC20_LOCK,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}
	packets, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*out = OLTReply{
		RawTX: packets,
	}
	return nil
}
