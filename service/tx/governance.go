package tx

import (
	"errors"

	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	gov "github.com/Oneledger/protocol/action/governance"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (s *Service) CreateProposal(args client.CreateProposalRequest, reply *client.CreateTxReply) error {
	proposalType := governance.NewProposalType(args.ProposalType)
	if proposalType == governance.ProposalTypeInvalid {
		return errors.New("invalid proposal type")
	}

	createProposal := gov.CreateProposal{
		ProposalID:      governance.ProposalID(args.ProposalID),
		ProposalType:    proposalType,
		Description:     args.Description,
		Headline:        args.Headline,
		Proposer:        args.Proposer,
		InitialFunding:  args.InitialFunding,
		FundingGoal:     args.FundingGoal,
		FundingDeadline: args.FundingDeadline,
		VotingDeadline:  args.VotingDeadline,
		PassPercentage:  args.PassPercentage,
		ConfigUpdate:    args.ConfigUpdate,
	}

	data, err := createProposal.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.PROPOSAL_CREATE,
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

func (s *Service) FundProposal(args client.FundProposalRequest, reply *client.CreateTxReply) error {

	fundProposal := gov.FundProposal{
		ProposalId:    args.ProposalId,
		FunderAddress: args.FunderAddress,
		FundValue:     args.FundValue,
	}

	data, err := fundProposal.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.PROPOSAL_FUND,
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

func (s *Service) CancelProposal(args client.CancelProposalRequest, reply *client.CreateTxReply) error {
	cancelProposal := gov.CancelProposal{
		ProposalId: args.ProposalId,
		Proposer:   args.Proposer,
		Reason:     args.Reason,
	}

	data, err := cancelProposal.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.PROPOSAL_CANCEL,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet := tx.RawBytes()
	*reply = client.CreateTxReply{RawTx: packet}

	return nil
}

func (s *Service) WithdrawProposalFunds(args client.WithdrawFundsRequest, reply *client.CreateTxReply) error {

	withdrawProposal := gov.WithdrawFunds{
		ProposalID:    args.ProposalID,
		Funder:        args.Funder,
		WithdrawValue: args.WithdrawValue,
		Beneficiary:   args.Beneficiary,
	}

	data, err := withdrawProposal.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := &action.RawTx{
		Type: action.PROPOSAL_WITHDRAW_FUNDS,
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

func (s *Service) VoteProposal(args client.VoteProposalRequest, reply *client.VoteProposalReply) error {
	// this node address is voter
	hPub, err := s.nodeContext.ValidatorPubKey().GetHandler()
	if err != nil {
		s.logger.Error("error get public key handler", err)
		return codes.ErrLoadingNodeKey
	}
	address := hPub.Address()

	// get private key
	hPri, err := s.nodeContext.PrivVal().GetHandler()
	if err != nil {
		s.logger.Error("error get private key handler", err)
		return codes.ErrLoadingNodeKey
	}

	// prepare Tx struct
	if args.Opinion == governance.OPIN_UNKNOWN {
		return errors.New("invalid vote opinion")
	}

	voteProposal := gov.VoteProposal{
		ProposalID:       governance.ProposalID(args.ProposalId),
		Address:          args.Address,
		ValidatorAddress: address,
		Opinion:          args.Opinion,
	}

	data, err := voteProposal.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{
		Price: args.GasPrice,
		Gas:   args.Gas,
	}

	tx := action.RawTx{
		Type: action.PROPOSAL_VOTE,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	// validator signs Tx
	rawData := tx.RawBytes()
	pubkey := hPri.PubKey()
	signed, _ := hPri.Sign(rawData)

	// reply
	signature := action.Signature{Signed: signed, Signer: pubkey}
	*reply = client.VoteProposalReply{RawTx: rawData, Signature: signature}

	return nil
}
