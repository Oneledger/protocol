package tx

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/network_delegation"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/google/uuid"
)

func (s *Service) NetUndelegate(args client.NetUndelegateRequest, reply *client.CreateTxReply) error {

	undelegate := network_delegation.Undelegate{
		Delegator: args.Delegator,
		Amount:    args.Amount,
	}

	data, err := undelegate.Marshal()
	if err != nil {
		return err
	}

	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}
	uuidNew, _ := uuid.NewUUID()
	tx := &action.RawTx{
		Type: action.NETWORK_DELEGATION_UNDELEGATE,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{RawTx: packet}

	return nil
}
