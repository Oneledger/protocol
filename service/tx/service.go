package tx

import (
	"fmt"

	"github.com/Oneledger/protocol/action"
	ev "github.com/Oneledger/protocol/action/evidence"
	action_sc "github.com/Oneledger/protocol/action/smart_contract"
	"github.com/Oneledger/protocol/action/staking"
	"github.com/Oneledger/protocol/action/transfer"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/google/uuid"
)

func Name() string {
	return "tx"
}

type Service struct {
	balances      *balance.Store
	router        action.Router
	accounts      accounts.Wallet
	validators    *identity.ValidatorStore
	govern        *governance.Store
	delegators    *delegation.DelegationStore
	evidenceStore *evidence.EvidenceStore
	feeOpt        *fees.FeeOption
	logger        *log.Logger
	nodeContext   node.Context
}

func NewService(
	balances *balance.Store,
	router action.Router,
	accounts accounts.Wallet,
	validators *identity.ValidatorStore,
	govern *governance.Store,
	delegators *delegation.DelegationStore,
	evidenceStore *evidence.EvidenceStore,
	feeOpt *fees.FeeOption,
	nodeCtx node.Context,
	logger *log.Logger,
) *Service {
	return &Service{
		balances:      balances,
		router:        router,
		nodeContext:   nodeCtx,
		accounts:      accounts,
		validators:    validators,
		govern:        govern,
		delegators:    delegators,
		evidenceStore: evidenceStore,
		feeOpt:        feeOpt,
		logger:        logger,
	}
}

func getSendTxContext(args client.SendTxRequest) (data []byte, t action.Type, err error) {
	if len(args.Data) != 0 {
		// means that we execute a contract method
		msg := action_sc.Execute{
			From:   keys.Address(args.From),
			Amount: args.Amount,
			Data:   args.Data,
			Nonce:  args.Nonce,
		}
		if len(args.To.Bytes()) != 0 {
			to := keys.Address(args.To)
			msg.To = &to
		}
		data, err = msg.Marshal()
		t = action.SC_EXECUTE
	} else {
		// just a regular send
		msg := transfer.Send{
			From:   keys.Address(args.From),
			To:     keys.Address(args.To),
			Amount: args.Amount,
		}
		data, err = msg.Marshal()
		t = action.SEND
	}
	return
}

// SendTx exists for maintaining backwards compatibility with existing olclient implementations. It returns
// a signed transaction
// TODO: deprecate this
func (svc *Service) SendTx(args client.SendTxRequest, reply *client.CreateTxReply) error {
	data, t, err := getSendTxContext(args)
	if err != nil {
		svc.logger.Error("error in serializing send object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
	tx := action.RawTx{
		Type: t,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	if _, err := svc.accounts.GetAccount(args.From); err != nil {
		return accounts.ErrGetAccountByAddress
	}

	pubKey, signed, err := svc.accounts.SignWithAddress(tx.RawBytes(), keys.Address(args.From))
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
	data, t, err := getSendTxContext(args)
	if err != nil {
		svc.logger.Error("error in serializing send object", err)
		return codes.ErrSerialization
	}

	uuidNew, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	fee := action.Fee{args.GasPrice, args.Gas}
	tx := &action.RawTx{
		Type: t,
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

func (svc *Service) CreateRawSendPool(args client.SendPoolTxRequest, reply *client.CreateTxReply) error {
	sendpool := transfer.SendPool{
		From:     args.From,
		PoolName: args.PoolName,
		Amount:   args.Amount,
	}
	data, err := sendpool.Marshal()
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

func (svc *Service) Vote(args client.VoteRequest, reply *client.VoteReply) error {
	if svc.evidenceStore.IsFrozenValidator(args.Address) {
		svc.logger.Error("got allegation validator", args.Address)
		return evidence.ErrFrozenValidator
	}

	validator, err := svc.validators.Get(args.Address)
	if err != nil {
		svc.logger.Errorf("validator for address %s not found\n", args.Address)
		return codes.ErrBadOwner
	}

	svc.logger.Detailf("Vote transaction for address - %s\n", validator.Address)

	allegation := ev.AllegationVote{
		Address:   validator.Address,
		RequestID: args.RequestID,
		Choice:    args.Choice,
	}

	data, err := allegation.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing vote object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := svc.feeOpt.MinFee()

	tx := action.RawTx{
		Type: action.ALLEGATION_VOTE,
		Data: data,
		Fee: action.Fee{
			Price: action.Amount{
				Currency: "OLT",
				Value:    *feeAmount.Amount,
			},
			Gas: 100000,
		},
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		svc.logger.Error("error in serializing send transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.VoteReply{RawTx: packet}

	return nil
}

func (svc *Service) Allegation(args client.AllegationRequest, reply *client.AllegationReply) error {
	mv, err := svc.validators.Get(args.MaliciousAddress)
	if err != nil {
		svc.logger.Errorf("malicious validator for address %s not found\n", args.MaliciousAddress)
		return codes.ErrBadOwner
	}

	if svc.evidenceStore.IsFrozenValidator(mv.Address) {
		svc.logger.Error("validator already frozen", mv.Address)
		return evidence.ErrFrozenValidator
	}

	validator, err := svc.validators.Get(args.Address)
	if err != nil {
		svc.logger.Errorf("validator for address %s not found\n", args.Address)
		return codes.ErrBadOwner
	}

	if svc.evidenceStore.IsFrozenValidator(validator.Address) {
		svc.logger.Error("validator could no perform allegation while frosted", validator.Address)
		return evidence.ErrFrozenValidator
	}

	requestID, err := svc.evidenceStore.GenerateRequestID()
	if err != nil {
		svc.logger.Errorf("failed to generate request ID\n")
		return codes.ErrSerialization
	}

	allegation := ev.Allegation{
		RequestID:        requestID,
		ValidatorAddress: validator.Address,
		MaliciousAddress: mv.Address,
		BlockHeight:      args.BlockHeight,
		ProofMsg:         args.ProofMsg,
	}

	data, err := allegation.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing allegation object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := svc.feeOpt.MinFee()

	tx := action.RawTx{
		Type: action.ALLEGATION,
		Data: data,
		Fee: action.Fee{
			Price: action.Amount{
				Currency: "OLT",
				Value:    *feeAmount.Amount,
			},
			Gas: 100000,
		},
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		svc.logger.Error("error in serializing send transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.AllegationReply{RawTx: packet}

	return nil
}

func (svc *Service) Release(args client.ReleaseRequest, reply *client.ReleaseReply) error {
	validator, err := svc.validators.Get(args.Address)
	svc.validators.Iterate(func(addr keys.Address, validator *identity.Validator) bool {
		fmt.Println("Validator :", validator.Address.String())
		return false
	})
	if err != nil {
		svc.logger.Errorf("validator for address %s not found\n", args.Address)
		return codes.ErrBadOwner
	}

	if !svc.evidenceStore.IsFrozenValidator(validator.Address) {
		svc.logger.Error("validator is not frozen", validator.Address)
		return evidence.ErrNonFrozenValidator
	}

	release := ev.Release{
		ValidatorAddress: validator.Address,
	}

	data, err := release.Marshal()
	if err != nil {
		svc.logger.Error("error in serializing release object", err)
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	feeAmount := svc.feeOpt.MinFee()

	tx := action.RawTx{
		Type: action.RELEASE,
		Data: data,
		Fee: action.Fee{
			Price: action.Amount{
				Currency: "OLT",
				Value:    *feeAmount.Amount,
			},
			Gas: 100000,
		},
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		svc.logger.Error("error in serializing send transaction", err)
		return codes.ErrSerialization
	}

	*reply = client.ReleaseReply{RawTx: packet}

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

	svc.logger.Infof("Validator - %s, delegator - %s, unstake amount - %+v\n",
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

	svc.logger.Infof("Validator - %s, delegator - %s, withdraw amount - %+v\n",
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
