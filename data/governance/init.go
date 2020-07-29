package governance

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
)

var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("governance")
}

const (
	//Proposal Types
	ProposalTypeInvalid      ProposalType = 0xEE
	ProposalTypeConfigUpdate ProposalType = 0x20
	ProposalTypeCodeChange   ProposalType = 0x21
	ProposalTypeGeneral      ProposalType = 0x22

	//Proposal Status
	ProposalStatusFunding   ProposalStatus = 0x23
	ProposalStatusVoting    ProposalStatus = 0x24
	ProposalStatusCompleted ProposalStatus = 0x25

	//Proposal Outcome
	ProposalOutcomeInProgress        ProposalOutcome = 0x26
	ProposalOutcomeInsufficientFunds ProposalOutcome = 0x27
	ProposalOutcomeInsufficientVotes ProposalOutcome = 0x28
	ProposalOutcomeCompletedNo       ProposalOutcome = 0x29
	ProposalOutcomeCancelled         ProposalOutcome = 0x30
	ProposalOutcomeCompletedYes      ProposalOutcome = 0x31

	//Proposal States
	ProposalStateInvalid        ProposalState = 0xEE
	ProposalStateActive         ProposalState = 0x32
	ProposalStatePassed         ProposalState = 0x33
	ProposalStateFailed         ProposalState = 0x34
	ProposalStateFinalized      ProposalState = 0x35
	ProposalStateFinalizeFailed ProposalState = 0x36

	//Vote Opinions
	OPIN_UNKNOWN  VoteOpinion = 0x0
	OPIN_POSITIVE VoteOpinion = 0x1
	OPIN_NEGATIVE VoteOpinion = 0x2
	OPIN_GIVEUP   VoteOpinion = 0x3

	//Vote Result
	VOTE_RESULT_PASSED VoteResult = 0x10
	VOTE_RESULT_FAILED VoteResult = 0x11
	VOTE_RESULT_TBD    VoteResult = 0x12

	//Error Codes
	errorSerialization   = "321"
	errorDeSerialization = "322"
	errorSettingRecord   = "323"
	errorGettingRecord   = "324"
	errorDeletingRecord  = "325"

	//Proposal ID length based on hash algorithm
	SHA256LENGTH int = 0x40
)

type GovProposal struct {
	Prop          Proposal        `json:"proposal"`
	ProposalVotes []*ProposalVote `json:"proposalVotes"`
	ProposalFunds []ProposalFund  `json:"proposalFunds"`
	State         ProposalState
}

type ProposalMasterStore struct {
	Proposal     *ProposalStore
	ProposalFund *ProposalFundStore
	ProposalVote *ProposalVoteStore
}

func (p *ProposalMasterStore) WithState(state *storage.State) *ProposalMasterStore {
	p.Proposal.WithState(state)
	p.ProposalFund.WithState(state)
	p.ProposalVote.WithState(state)
	return p
}

func (p *ProposalMasterStore) GetProposalVotes(id ProposalID) []*ProposalVote {
	_, votes, err := p.ProposalVote.GetVotesByID(id)
	if err != nil {
		return nil
	}
	return votes
}

func (p *ProposalMasterStore) GetProposalFunds(id ProposalID) []ProposalFund {
	return p.ProposalFund.GetFundsForProposalID(id, func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund {
		propFund := ProposalFund{
			Id:            proposalID,
			Address:       fundingAddr,
			FundingAmount: amt,
		}
		return propFund
	})
}

func (p *ProposalMasterStore) LoadProposals(proposals []GovProposal) error {
	for _, prop := range proposals {
		p.Proposal.WithPrefixType(prop.State)
		err := p.Proposal.Set(&prop.Prop)
		if err != nil {
			return err
		}
		for _, vote := range prop.ProposalVotes {
			err = p.ProposalVote.Setup(prop.Prop.ProposalID, NewProposalVote(vote.Validator, OPIN_UNKNOWN, vote.Power))
			if err != nil {
				return err
			}
			err = p.ProposalVote.Update(prop.Prop.ProposalID, vote)
			if err != nil {
				return err
			}
		}
		for _, fund := range prop.ProposalFunds {
			err = p.ProposalFund.AddFunds(fund.Id, fund.Address, fund.FundingAmount)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewProposalMasterStore(p *ProposalStore, pf *ProposalFundStore, pv *ProposalVoteStore) *ProposalMasterStore {
	return &ProposalMasterStore{
		Proposal:     p,
		ProposalFund: pf,
		ProposalVote: pv,
	}
}
