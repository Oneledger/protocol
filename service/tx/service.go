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
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
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
	validators  *identity.ValidatorStore
	govern      *governance.Store
	delegators  *delegation.DelegationStore
	feeOpt      *fees.FeeOption
	logger      *log.Logger
	nodeContext node.Context
}

func NewService(
	balances *balance.Store,
	router action.Router,
	accounts accounts.Wallet,
	validators *identity.ValidatorStore,
	govern *governance.Store,
	delegators *delegation.DelegationStore,
	feeOpt *fees.FeeOption,
	nodeCtx node.Context,
	logger *log.Logger,
) *Service {
	return &Service{
		balances:    balances,
		router:      router,
		nodeContext: nodeCtx,
		accounts:    accounts,
		validators:  validators,
		govern:      govern,
		delegators:  delegators,
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

func (svc *Service) Stake(args client.StakeRequest, reply *client.StakeReply) error {
	if len(args.Name) < 1 {
		args.Name = svc.nodeContext.NodeName
	}

	pubkey := svc.nodeContext.ValidatorPubKey()
	ecdsaPubKey := svc.nodeContext.ValidatorECDSAPubKey()
	handler, _ := pubkey.GetHandler()
	address := handler.Address()

	svc.logger.Infof("Validator - %s, delegator - %s, stake amount - %+v\n",
		address, args.Address, args.Amount,
	)

	apply := staking.Stake{
		ValidatorAddress:     address,
		StakeAddress:         args.Address,
		Stake:                action.Amount{Currency: "OLT", Value: args.Amount},
		NodeName:             args.Name,
		ValidatorPubKey:      pubkey,
		ValidatorECDSAPubKey: ecdsaPubKey,
	}

	data, err := apply.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing Stake object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := svc.feeOpt.MinFee()

	tx := action.RawTx{
		Type: action.STAKE,
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

	*reply = client.StakeReply{RawTx: packet}

	return nil
}

func (svc *Service) Unstake(args client.UnstakeRequest, reply *client.UnstakeReply) error {
	pubkey := svc.nodeContext.ValidatorPubKey()
	handler, _ := pubkey.GetHandler()
	address := handler.Address()

	svc.logger.Infof("Validator - %s, delegator - %s, stake amount - %+v\n",
		address, args.Address, args.Amount,
	)

	apply := staking.Unstake{
		ValidatorAddress: address,
		StakeAddress:     args.Address,
		Stake:            action.Amount{Currency: "OLT", Value: args.Amount},
	}

	data, err := apply.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing Unstake object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := svc.feeOpt.MinFee()

	tx := action.RawTx{
		Type: action.UNSTAKE,
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

	*reply = client.UnstakeReply{RawTx: packet}

	return nil
}

func (svc *Service) Withdraw(args client.WithdrawRequest, reply *client.WithdrawReply) error {
	pubkey := svc.nodeContext.ValidatorPubKey()
	handler, _ := pubkey.GetHandler()
	address := handler.Address()

	svc.logger.Infof("Validator - %s, delegator - %s, stake amount - %+v\n",
		address, args.Address, args.Amount,
	)

	apply := staking.Withdraw{
		ValidatorAddress: address,
		StakeAddress:     args.Address,
		Stake:            action.Amount{Currency: "OLT", Value: args.Amount},
	}

	data, err := apply.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing Withdraw object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := svc.feeOpt.MinFee()

	tx := action.RawTx{
		Type: action.WITHDRAW,
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

	*reply = client.WithdrawReply{RawTx: packet}

	return nil
}
