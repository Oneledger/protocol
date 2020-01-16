/*

 */

package tx

import (
	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/ons"
	"github.com/Oneledger/protocol/client"
	ons2 "github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (s *Service) ONS_CreateRawCreate(args client.ONSCreateRequest, reply *client.CreateTxReply) error {

	name := ons2.GetNameFromString(args.Name)
	domainCreate := ons.DomainCreate{
		Owner:       args.Owner,
		Beneficiary: args.Account,
		Name:        name,
		Uri:         args.Uri,
		BuyingPrice: args.BuyingPrice,
	}

	data, err := domainCreate.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
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

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawUpdate(args client.ONSUpdateRequest, reply *client.CreateTxReply) error {

	name := ons2.GetNameFromString(args.Name)
	domainUpdate := ons.DomainUpdate{
		Owner:       args.Owner,
		Beneficiary: args.Account,
		Name:        name,
		Active:      args.Active,
		Uri:         args.Uri,
	}
	data, err := domainUpdate.Marshal()
	if err != nil {
		s.logger.Error("error in serializing domain update object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
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

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawRenew(args client.ONSRenewRequest, reply *client.CreateTxReply) error {

	name := ons2.GetNameFromString(args.Name)
	renewDomain := ons.RenewDomain{
		Owner:       args.Owner,
		Name:        name,
		BuyingPrice: args.BuyingPrice,
	}
	data, err := renewDomain.Marshal()
	if err != nil {
		s.logger.Error("error in serializing domain update object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
	tx := &action.RawTx{
		Type: action.DOMAIN_RENEW,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		s.logger.Error("error in serializing domain update transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawSale(args client.ONSSaleRequest, reply *client.CreateTxReply) error {

	name := ons2.GetNameFromString(args.Name)
	domainSale := ons.DomainSale{
		Name:         name,
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
	fee := action.Fee{args.GasPrice, args.Gas}
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

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawBuy(args client.ONSPurchaseRequest, reply *client.CreateTxReply) error {

	name := ons2.GetNameFromString(args.Name)
	domainPurchase := ons.DomainPurchase{
		Name:     name,
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
	fee := action.Fee{args.GasPrice, args.Gas}
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

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawSend(args client.ONSSendRequest, reply *client.CreateTxReply) error {

	name := ons2.GetNameFromString(args.Name)
	domainSend := ons.DomainSend{
		DomainName: name,
		From:       args.From,
		Amount:     args.Amount,
	}
	data, err := domainSend.Marshal()
	if err != nil {
		s.logger.Error("error in serializing domain send object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
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

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) ONS_CreateRawDeleteSub(args client.ONSDeleteSubRequest, reply *client.CreateTxReply) error {

	name := ons2.GetNameFromString(args.Name)
	del := ons.DeleteSub{
		Name:  name,
		Owner: args.Owner,
	}
	data, err := del.Marshal()
	if err != nil {
		s.logger.Error("error in serializing domain delete object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
	tx := &action.RawTx{
		Type: action.DOMAIN_DELETE_SUB,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		s.logger.Error("error in serializing domain delete transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{
		RawTx: packet,
	}
	return nil
}
