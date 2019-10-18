/*

 */

package tx

import (
	"strings"

	codes "github.com/Oneledger/protocol/status_codes"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/ons"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	"github.com/google/uuid"
)

func (s *Service) ONS_CreateRawCreate(args client.ONSCreateRequest, reply *client.SendTxReply) error {

	domainCreate := ons.DomainCreate{
		Owner:   args.Owner,
		Account: args.Account,
		Name:    args.Name,
		Price:   args.Price,
	}
	data, err := domainCreate.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.RawTx{
		Type: action.DOMAIN_CREATE,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawUpdate(args client.ONSUpdateRequest, reply *client.SendTxReply) error {

	domainUpdate := ons.DomainUpdate{
		Owner:   args.Owner,
		Account: args.Account,
		Name:    args.Name,
		Active:  args.Active,
	}
	data, err := domainUpdate.Marshal()
	if err != nil {
		s.logger.Error("error in serializing domain update object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.RawTx{
		Type: action.DOMAIN_UPDATE,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		s.logger.Error("error in serializing domain update transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawSale(args client.ONSSaleRequest, reply *client.SendTxReply) error {

	domainSale := ons.DomainSale{
		DomainName:   args.Name,
		OwnerAddress: args.OwnerAddress,
		Price:        args.Price,
		CancelSale:   args.CancelSale,
	}
	data, err := domainSale.Marshal()
	if err != nil {
		s.logger.Error("error in serializing domain sale object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.RawTx{
		Type: action.DOMAIN_SELL,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		s.logger.Error("error in serializing domain sale transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawBuy(args client.ONSPurchaseRequest, reply *client.SendTxReply) error {

	domainPurchase := ons.DomainPurchase{
		Name:     strings.ToLower(args.Name),
		Buyer:    args.Buyer,
		Account:  args.Account,
		Offering: args.Offering,
	}
	data, err := domainPurchase.Marshal()
	if err != nil {
		s.logger.Error("error in serializing domain purchase object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.RawTx{
		Type: action.DOMAIN_PURCHASE,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		s.logger.Error("error in serializing domain purchase transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawSend(args client.ONSSendRequest, reply *client.SendTxReply) error {

	domainSend := ons.DomainSend{
		DomainName: args.Name,
		From:       args.From,
		Amount:     args.Amount,
	}
	data, err := domainSend.Marshal()
	if err != nil {
		s.logger.Error("error in serializing domain send object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.RawTx{
		Type: action.DOMAIN_SEND,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		s.logger.Error("error in serializing domain send transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}
