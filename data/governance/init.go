package governance

import (
	"os"

	"github.com/Oneledger/protocol/log"
)

var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("governance")
}

const (
	//Proposal Types
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
	ProposalStateError  ProposalState = 0xEE
	ProposalStateActive ProposalState = 0x31
	ProposalStatePassed ProposalState = 0x32
	ProposalStateFailed ProposalState = 0x33

	//Vote Opinions
	OPIN_UNKNOWN  VoteOpinion = 0
	OPIN_POSITIVE VoteOpinion = 1
	OPIN_NEGATIVE VoteOpinion = 2
	OPIN_GIVEUP   VoteOpinion = 3

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
}

func NewProposalMasterStore(p *ProposalStore, pf *ProposalFundStore) *ProposalMasterStore {
	return &ProposalMasterStore{
		Proposal:     p,
		ProposalFund: pf,
	}
}
