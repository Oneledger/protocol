package ethereum

import (
	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/serialize"
)

func (svc *Service) CreateRawExtRedeem(req RedeemRequest, out *OLTLockReply) error {

	redeem := eth.Redeem{
		Owner:  req.userOLTaddress,
		To:     req.userETHaddress,
		ETHTxn: req.ETHTxn,
	}

	data, err := redeem.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{Price: req.Fee, Gas: req.Gas}
	tx := &action.RawTx{
		Type: action.ETH_REDEEM,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return action.ErrUnserializable
	}
	*out = OLTLockReply{
		RawTX: packet,
	}
	return nil
}