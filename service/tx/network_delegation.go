package tx

import (
	"github.com/Oneledger/protocol/action"
	nwd "github.com/Oneledger/protocol/action/network_delegation"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/google/uuid"
)

func (s *Service) AddNetworkDelegation(args client.NetworkDelegateRequest, reply *client.CreateTxReply) error {

	networkDelegation := nwd.AddNetworkDelegation{
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
		Type: action.ADD_NETWORK_DELEGATION,
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

func (s *Service) NetworkUndelegate(args client.NetUndelegateRequest, reply *client.CreateTxReply) error {

	undelegate := nwd.Undelegate{
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
		Type: action.NETWORK_UNDELEGATE,
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
