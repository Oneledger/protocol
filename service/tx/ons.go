/*

 */

package tx

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/ons"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/rpc"
	"github.com/Oneledger/protocol/serialize"
	"github.com/google/uuid"
)

func (s *Service) Ons_CreateRawCreate(args client.ONSCreateRequest, reply *client.SendTxReply) error {

	domainCreate := ons.DomainCreate{
		Owner:   args.Owner,
		Account: args.Account,
		Name:    args.Name,
		Price:   args.Price,
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.BaseTx{
		Data: domainCreate,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return rpc.NewError(rpc.CodeInternalError, "err while network serialization:"+err.Error())
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) Ons_CreateRawUpdate(args client.ONSUpdateRequest, reply *client.SendTxReply) error {

	domainUpdate := ons.DomainUpdate{
		Owner:   args.Owner,
		Account: args.Account,
		Name:    args.Name,
		Active:  args.Active,
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.BaseTx{
		Data: domainUpdate,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return rpc.NewError(rpc.CodeInternalError, "err while network serialization:"+err.Error())
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) Ons_CreateRawSale(args client.ONSSaleRequest, reply *client.SendTxReply) error {

	domainSale := ons.DomainSale{
		DomainName:   args.Name,
		OwnerAddress: args.OwnerAddress,
		Price:        args.Price,
		CancelSale:   args.CancelSale,
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.BaseTx{
		Data: domainSale,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return rpc.NewError(rpc.CodeInternalError, "err while network serialization:"+err.Error())
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) Ons_CreateRawBuy(args client.ONSPurchaseRequest, reply *client.SendTxReply) error {

	domainPurchase := ons.DomainPurchase{
		Name:     args.Name,
		Buyer:    args.Buyer,
		Account:  args.Account,
		Offering: args.Offering,
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.BaseTx{
		Data: domainPurchase,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return rpc.NewError(rpc.CodeInternalError, "err while network serialization:"+err.Error())
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) Ons_CreateRawSend(args client.ONSSendRequest, reply *client.SendTxReply) error {

	domainSend := ons.DomainSend{
		DomainName: args.Name,
		From:       args.From,
		Amount:     args.Amount,
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.BaseTx{
		Data: domainSend,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return rpc.NewError(rpc.CodeInternalError, "err while network serialization"+err.Error())
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}
