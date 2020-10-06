package tx

import (
	"github.com/Oneledger/protocol/action"
	nwd "github.com/Oneledger/protocol/action/network_delegation"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/google/uuid"
)

func (s *Service) NetworkDelegate(args client.NetworkDelegateRequest, reply *client.CreateTxReply) error {

	networkDelegation := nwd.NetworkDelegate{
		DelegationAddress: args.DelegationAddress,
		Amount:            args.Amount,
	}

	data, err := networkDelegation.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.NETWORKDELEGATE,
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

func (s *Service) WithDrawNetworkDelegation(args client.NetworkDelegateRequest, reply *client.CreateTxReply) error {

	withdrawNetworkDelegation := nwd.WithdrawNetworkDelegation{
		DelegationAddress: args.DelegationAddress,
		Amount:            args.Amount,
	}

	data, err := withdrawNetworkDelegation.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.WITHDRAW_NETWORK_DELEGATION,
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
