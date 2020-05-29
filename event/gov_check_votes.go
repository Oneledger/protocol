package event

import (
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
	"strconv"
)

var _ jobs.Job = &JobGovCheckVotes{}

type JobGovCheckVotes struct {
	ProposalID governance.ProposalID
	JobID      string
	RetryCount int
	Status     jobs.Status
}

func NewGovCheckVotes(proposalID governance.ProposalID, status governance.ProposalStatus) *JobGovCheckVotes {
	return &JobGovCheckVotes{
		ProposalID: proposalID,
		JobID:      proposalID.String() + storage.DB_PREFIX + strconv.Itoa(int(status)),
		RetryCount: 0,
		Status:     0,
	}
}

func (j JobGovCheckVotes) DoMyJob(ctx interface{}) {
	govCtx, _ := ctx.(*JobsContext)
	proposalMaster := govCtx.ProposalMaster

	//Get Proposal
	proposal, err := proposalMaster.Proposal.Get(j.ProposalID)
	if err != nil {
		j.Status = jobs.Failed
		govCtx.Logger.Error("gov_check_votes: failed to get proposal")
		return
	}

	//Check if proposal is in voting state
	if proposal.Status != governance.ProposalStatusVoting {
		j.Status = jobs.Failed
		govCtx.Logger.Error("gov_check_votes: proposal is not in a voting state")
		return
	}

	//Check number of votes
	passed, err := proposalMaster.ProposalVote.ResultSoFar(j.ProposalID, proposal.PassPercentage)
	if err != nil {
		j.Status = jobs.Failed
		govCtx.Logger.Error(errors.Wrap(err, "gov_check_votes:"))
		return
	}
	if passed == governance.VOTE_RESULT_PASSED {
		j.Status = jobs.Failed
		govCtx.Logger.Error("gov_check_votes: proposal has already been passed")
		return
	}

	//Check vote expiry
	if proposalMaster.Proposal.GetState().Version() <= proposal.VotingDeadline {
		j.Status = jobs.Failed
		govCtx.Logger.Error("gov_check_votes: voting period is not over yet")
		return
	}

	//Create internal transaction and broadcast
	err = BroadcastGovExpireVotesTx(govCtx, j.ProposalID, j.JobID)
	if err != nil {
		j.Status = jobs.Failed
		govCtx.Logger.Error("gov_check_votes: error broadcasting" + err.Error())
		return
	}

	j.Status = jobs.Completed
}

func (j JobGovCheckVotes) IsDone() bool {
	return j.Status == jobs.Completed
}

func (j JobGovCheckVotes) IsFailed() bool {
	return j.Status == jobs.Failed
}

func (j JobGovCheckVotes) GetType() string {
	return JobTypeGOVCheckVotes
}

func (j JobGovCheckVotes) GetJobID() string {
	return j.JobID
}
