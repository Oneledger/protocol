package tx

import (
	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/staking"
	"github.com/Oneledger/protocol/action/transfer"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func Name() string {
	return "tx"
}

type Service struct {
	balances    *balance.Store
	router      action.Router
	accounts    accounts.Wallet
	feeOpt      *fees.FeeOption
	logger      *log.Logger
	nodeContext node.Context
}

func NewService(
	balances *balance.Store,
	router action.Router,
	accounts accounts.Wallet,
	feeOpt *fees.FeeOption,
	nodeCtx node.Context,
	logger *log.Logger,
) *Service {
	return &Service{
		balances:    balances,
		router:      router,
		nodeContext: nodeCtx,
		accounts:    accounts,
		feeOpt:      feeOpt,
		logger:      logger,
	}
}

// SendTx exists for maintaining backwards compatibility with existing olclient implementations. It returns
// a signed transaction
// TODO: deprecate this
func (svc *Service) SendTx(args client.SendTxRequest, reply *client.CreateTxReply) error {
	send := transfer.Send{
		From:   keys.Address(args.From),
		To:     keys.Address(args.To),
		Amount: args.Amount,
	}
	data, err := send.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing send object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
	tx := action.RawTx{
		Type: action.SEND,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	if _, err := svc.accounts.GetAccount(args.From); err != nil {
		return accounts.ErrGetAccountByAddress
	}

	pubKey, signed, err := svc.accounts.SignWithAddress(tx.RawBytes(), send.From)
	if err != nil {
		return err
	}
	signatures := []action.Signature{{pubKey, signed}}
	signedTx := &action.SignedTx{
		RawTx:      tx,
		Signatures: signatures,
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (svc *Service) SendPoolTx(args client.SendPoolTxRequest, reply *client.CreateTxReply) error {
	send := transfer.SendPool{
		From:     args.From,
		PoolName: args.PoolName,
		Amount:   args.Amount,
	}
	data, err := send.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing sendPool object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
	tx := action.RawTx{
		Type: action.SENDPOOL,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	if _, err := svc.accounts.GetAccount(args.From); err != nil {
		return accounts.ErrGetAccountByAddress
	}

	pubKey, signed, err := svc.accounts.SignWithAddress(tx.RawBytes(), send.From)
	if err != nil {
		return err
	}
	signatures := []action.Signature{{pubKey, signed}}
	signedTx := &action.SignedTx{
		RawTx:      tx,
		Signatures: signatures,
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (svc *Service) CreateRawSend(args client.SendTxRequest, reply *client.CreateTxReply) error {
	send := transfer.Send{
		From:   args.From,
		To:     args.To,
		Amount: args.Amount,
	}
	data, err := send.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing send object", err)
		return codes.ErrSerialization
	}

	uuidNew, err := uuid.NewUUID()

	fee := action.Fee{args.GasPrice, args.Gas}
	tx := &action.RawTx{
		Type: action.SEND,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		svc.logger.Error("error in serializing send transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{
		RawTx: packet,
	}

	return nil
}

func (svc *Service) ApplyValidator(args client.ApplyValidatorRequest, reply *client.ApplyValidatorReply) error {
	if len(args.Name) < 1 {
		args.Name = svc.nodeContext.NodeName
	}

	if len(args.Address) < 1 {
		handler, err := svc.nodeContext.PubKey().GetHandler()
		if err != nil {
			return err
		}
		args.Address = handler.Address()
	}

	pubkey := &keys.PublicKey{keys.GetAlgorithmFromTmKeyName(args.TmPubKeyType), args.TmPubKey}
	if len(args.TmPubKey) < 1 {
		*pubkey = svc.nodeContext.ValidatorPubKey()
	}

	ecdsaPubKey := svc.nodeContext.ValidatorECDSAPubKey()

	handler, err := pubkey.GetHandler()
	if err != nil {

		return err
	}

	addr := handler.Address()
	apply := staking.ApplyValidator{
		StakeAddress:         args.Address,
		Stake:                action.Amount{Currency: "VT", Value: args.Amount},
		NodeName:             args.Name,
		ValidatorAddress:     addr,
		ValidatorPubKey:      *pubkey,
		ValidatorECDSAPubKey: ecdsaPubKey,
		Purge:                args.Purge,
	}

	data, err := apply.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing applyValidator object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := svc.feeOpt.MinFee()

	tx := action.RawTx{
		Type: action.APPLYVALIDATOR,
		Data: data,
		Fee:  action.Fee{action.Amount{Currency: "OLT", Value: *feeAmount.Amount}, 100000},
		Memo: uuidNew.String(),
	}
	rawData := tx.RawBytes()
	h, err := svc.nodeContext.PrivVal().GetHandler()
	if err != nil {
		svc.logger.Error("error get validator handler", err)
		return codes.ErrLoadingNodeKey
	}
	vpubkey := h.PubKey()
	vsinged, err := h.Sign(rawData)

	signatures := []action.Signature{{vpubkey, vsinged}}
	signedTx := &action.SignedTx{
		RawTx:      tx,
		Signatures: signatures,
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if err != nil {
		svc.logger.Error("error in serializing signed transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.ApplyValidatorReply{RawTx: packet}

	return nil
}

func (svc *Service) WithdrawReward(args client.WithdrawRewardRequest, reply *client.WithdrawRewardReply) error {

	if len(args.To) < 1 {
		args.To = args.From
	}

	withdraw := staking.Withdraw{
		From: args.From,
		To:   args.To,
	}

	data, err := withdraw.Marshal()
	if err != nil {
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()

	tx := action.RawTx{
		Type: action.WITHDRAW,
		Data: data,
		Fee: action.Fee{
			Price: args.GasPrice,
			Gas:   args.Gas,
		},
		Memo: uuidNew.String(),
	}

	packet := tx.RawBytes()

	*reply = client.WithdrawRewardReply{RawTx: packet}
	return nil
}

func (svc Service) PurgeValidator(req client.PurgeValidatorRequest, reply *client.PurgeValidatorReply) error {

	if req.Validator.Err() != nil || req.Admin.Err() != nil {
		return codes.ErrBadAddress
	}

	purge := staking.Purge{
		AdminAddress:     req.Admin,
		ValidatorAddress: req.Validator,
	}
	data, err := purge.Marshal()
	if err != nil {
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := svc.feeOpt.MinFee()
	tx := action.RawTx{
		Type: purge.Type(),
		Data: data,
		Fee: action.Fee{
			Price: action.Amount{Currency: "OLT", Value: *feeAmount.Amount},
			Gas:   80000,
		},
		Memo: uuidNew.String(),
	}

	packet := tx.RawBytes()
	*reply = client.PurgeValidatorReply{
		RawTx: packet,
	}
	return nil
}
