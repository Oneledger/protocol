package tx

import (
	"errors"
	"strings"

	"time"

	"github.com/Oneledger/protocol/action"
	gov "github.com/Oneledger/protocol/action/governance"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/google/uuid"
)

func translateProposalType(propType string) governance.ProposalType {
	switch propType {
	case "codeChange":
		return governance.ProposalTypeCodeChange
	case "configUpdate":
		return governance.ProposalTypeConfigUpdate
	case "general":
		return governance.ProposalTypeGeneral
	default:
		return governance.ProposalTypeError
	}
}

func (s *Service) CreateProposal(args client.CreateProposalRequest, reply *client.CreateTxReply) error {
	proposalType := translateProposalType(args.ProposalType)
	if proposalType == governance.ProposalTypeError {
		return errors.New("invalid proposal type")
	}

	uniqueStr := args.ProposalID + time.Now().String()

	createProposal := gov.CreateProposal{
		ProposalID:     governance.ProposalID(uniqueStr),
		ProposalType:   proposalType,
		Description:    args.Description,
		Proposer:       args.Proposer,
		InitialFunding: args.InitialFunding,
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

func translateVoteOpinion(opin string) governance.VoteOpinion {
	opin = strings.ToUpper(opin)
	switch opin {
	case "YES":
		return governance.OPIN_POSITIVE
	case "NO":
		return governance.OPIN_NEGATIVE
	case "GIVEUP":
		return governance.OPIN_GIVEUP
	default:
		return governance.OPIN_UNKNOWN
	}
}

func (s *Service) CreateVote(args client.CreateVoteRequest, reply *client.CreateVoteReply) error {
	// this node address is voter
	hPub, err := s.nodeContext.PubKey().GetHandler()
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
	opin := translateVoteOpinion(args.Opinion)
	if opin == governance.OPIN_UNKNOWN {
		return errors.New("invalid vote opinion")
	}

	voteProposal := gov.VoteProposal{
		ProposalID: governance.IDFromString(args.ProposalID),
		Address:    address,
		Opinion:    opin,
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

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	// validator signs Tx
	rawData := tx.RawBytes()
	pubkey := hPri.PubKey()
	signed, _ := hPri.Sign(rawData)

	// reply
	signature := action.Signature{Signed: signed, Signer: pubkey}
	*reply = client.CreateVoteReply{RawTx: packet, Signature: signature}

	return nil
}
