package governance

import (
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
	ProposalTypeError        ProposalType = 0xEE
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
	ProposalOutcomeCancelled         ProposalOutcome = 0x29
	ProposalOutcomeCompleted         ProposalOutcome = 0x30

	//Proposal States
	ProposalStateError          ProposalState = 0xEE
	ProposalStateActive         ProposalState = 0x31
	ProposalStatePassed         ProposalState = 0x32
	ProposalStateFailed         ProposalState = 0x33
	ProposalStateFinalized      ProposalState = 0x34
	ProposalStateFinalizeFailed ProposalState = 0x35

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
)

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

func NewProposalMasterStore(p *ProposalStore, pf *ProposalFundStore, pv *ProposalVoteStore) *ProposalMasterStore {
	return &ProposalMasterStore{
		Proposal:     p,
		ProposalFund: pf,
		ProposalVote: pv,
	}
}
