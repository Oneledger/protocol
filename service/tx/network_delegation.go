package tx

import (
	"github.com/Oneledger/protocol/action"
	netwkDeleg "github.com/Oneledger/protocol/action/network_delegation"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"

	"github.com/google/uuid"
)

func (s *Service) WithdrawDelegRewards(args client.WithdrawDelegRewardsRequest, reply *client.CreateTxReply) error {
	withdraw := netwkDeleg.Withdraw{
		Delegator: args.Delegator,
		Amount:    action.Amount{Currency: "OLT", Value: args.Amount},
	}

	data, err := withdraw.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := s.feeOpt.MinFee()
	tx := &action.RawTx{
		Type: action.WITHDRAW_DELEG_REWARD,
		Data: data,
		Fee: action.Fee{
			Price: action.Amount{Currency: "OLT", Value: *feeAmount.Amount},
			Gas:   80000,
		},
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{RawTx: packet}
	return nil
}
