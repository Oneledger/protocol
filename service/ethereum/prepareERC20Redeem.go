package ethereum

import (
	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (svc *Service) CreateRawExtERC20Redeem(req RedeemRequest, out *OLTReply) error {

	redeemERC20 := eth.ERC20Redeem{
		Owner:  req.UserOLTaddress,
		To:     req.UserETHaddress,
		ETHTxn: req.ETHTxn,
	}

	data, err := redeemERC20.Marshal()
	if err != nil {
		svc.logger.Error(codes.ErrUnmarshaling.ErrorMsg())
		return codes.ErrUnmarshaling
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{Price: req.Fee, Gas: req.Gas}
	tx := &action.RawTx{
		Type: action.ERC20_REDEEM,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return action.ErrUnserializable
	}
	*out = OLTReply{
		RawTX: packet,
	}
	return nil
}