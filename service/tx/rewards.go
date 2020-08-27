package tx

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/rewards"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"

	"github.com/google/uuid"
)

func (s *Service) WithdrawRewards(args client.WithdrawRewardsRequest, reply *client.WithdrawRewardsReply) error {
	// Get address of local validator , assuming withdraw rewards will be called by validator in his own node
	hPub, err := s.nodeContext.ValidatorPubKey().GetHandler()
	if err != nil {
		s.logger.Error("error get public key handler", err)
		return codes.ErrLoadingNodeKey
	}
	address := hPub.Address()
	// create Withdrawal request for Validator
	withdrawRewards := rewards.Withdraw{
		ValidatorAddress: address,
		SignerAddress:    args.ValidatorAddress,
		WithdrawAmount:   action.Amount{Currency: "OLT", Value: args.WithdrawAmount},
	}

	data, err := withdrawRewards.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := s.feeOpt.MinFee()
	tx := &action.RawTx{
		Type: action.WITHDRAW_REWARD,
		Data: data,
		Fee:  action.Fee{action.Amount{Currency: "OLT", Value: *feeAmount.Amount}, 100000},
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.WithdrawRewardsReply{RawTx: packet}

	return nil
}
