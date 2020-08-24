package tx

import (
	"github.com/Oneledger/protocol/action/bidding"

	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (s *Service) CreateBid(args client.CreateBidRequest, reply *client.CreateTxReply) error {
	createBid := bidding.CreateBid{
		BidConvId: args.BidConvId,
		AssetOwner: args.AssetOwner,
		Asset: args.Asset,
		AssetType: args.AssetType,
		Bidder: args.Bidder,
		Amount: args.Amount,
		Deadline: args.Deadline,
	}

	data, err := createBid.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.BID_CREATE,
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

func (s *Service) CounterOffer(args client.CounterOfferRequest, reply *client.CreateTxReply) error {

	counterOffer := bidding.CounterOffer{
		BidConvId: args.BidConvId,
		AssetOwner: args.AssetOwner,
		Amount: args.Amount,
	}

	data, err := counterOffer.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.BID_CONTER_OFFER,
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

func (s *Service) CancelBid(args client.CancelBidRequest, reply *client.CreateTxReply) error {
	cancelBid := bidding.CancelBid{
		BidConvId: args.BidConvId,
		Bidder: args.Bidder,
	}

	data, err := cancelBid.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.BID_CANCEL,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet := tx.RawBytes()
	*reply = client.CreateTxReply{RawTx: packet}

	return nil
}

func (s *Service) OwnerDecision(args client.OwnerDecisionRequest, reply *client.CreateTxReply) error {
	ownerDecision := bidding.OwnerDecision{
		BidConvId: args.BidConvId,
		Owner: args.Owner,
		Decision: args.Decision,
	}

	data, err := ownerDecision.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.BID_OWNER_DECISION,
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

func (s *Service) BidderDecision(args client.BidderDecisionRequest, reply *client.CreateTxReply) error {
	bidderDecision := bidding.BidderDecision{
		BidConvId: args.BidConvId,
		Bidder: args.Bidder,
		Decision: args.Decision,
	}

	data, err := bidderDecision.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.BID_BIDDER_DECISION,
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
