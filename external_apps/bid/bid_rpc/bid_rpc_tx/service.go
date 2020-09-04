package bid_rpc_tx

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/external_apps/bid/bid_action"
	"github.com/Oneledger/protocol/external_apps/bid/bid_rpc"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/google/uuid"
)

func Name() string {
	return "bid_tx"
}

type Service struct {
	balances    *balance.Store
	router      action.Router
	logger      *log.Logger
}

func NewService(
	balances *balance.Store,
	logger *log.Logger,
) *Service {
	return &Service{
		balances:    balances,
		logger:      logger,
	}
}

func (s *Service) CreateBid(args bid_rpc.CreateBidRequest, reply *client.CreateTxReply) error {
	createBid := bid_action.CreateBid{
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
		Type: bid_action.BID_CREATE,
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

func (s *Service) CounterOffer(args bid_rpc.CounterOfferRequest, reply *client.CreateTxReply) error {

	counterOffer := bid_action.CounterOffer{
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
		Type: bid_action.BID_CONTER_OFFER,
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

func (s *Service) CancelBid(args bid_rpc.CancelBidRequest, reply *client.CreateTxReply) error {
	cancelBid := bid_action.CancelBid{
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
		Type: bid_action.BID_CANCEL,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet := tx.RawBytes()
	*reply = client.CreateTxReply{RawTx: packet}

	return nil
}

func (s *Service) OwnerDecision(args bid_rpc.OwnerDecisionRequest, reply *client.CreateTxReply) error {
	ownerDecision := bid_action.OwnerDecision{
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
		Type: bid_action.BID_OWNER_DECISION,
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

func (s *Service) BidderDecision(args bid_rpc.BidderDecisionRequest, reply *client.CreateTxReply) error {
	bidderDecision := bid_action.BidderDecision{
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
		Type: bid_action.BID_BIDDER_DECISION,
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