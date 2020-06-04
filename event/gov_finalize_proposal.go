package event

import (
	"strconv"

	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"
)

var _ jobs.Job = &JobGovfinalizeProposal{}

type JobGovfinalizeProposal struct {
	ProposalID governance.ProposalID
	JobID      string
	RetryCount int
	Status     jobs.Status
}

func NewGovFinalizeProposalJob(proposalID governance.ProposalID, status governance.ProposalStatus) *JobGovfinalizeProposal {
	return &JobGovfinalizeProposal{
		ProposalID: proposalID,
		JobID:      proposalID.String() + storage.DB_PREFIX + strconv.Itoa(int(status)),
		RetryCount: 0,
		Status:     0,
	}
}

func (j JobGovfinalizeProposal) DoMyJob(ctx interface{}) {
	govCtx, _ := ctx.(*JobsContext)
	BroadcastGovFinalizeVotesTx(govCtx, j.ProposalID, j.JobID)
	j.Status = jobs.Completed
}

func (j JobGovfinalizeProposal) IsDone() bool {
	return j.IsDone()
}

func (j JobGovfinalizeProposal) IsFailed() bool {
	return j.IsFailed()
}

func (j JobGovfinalizeProposal) GetType() string {
	return JobTypeGOVFinalizeProposal
}

func (j JobGovfinalizeProposal) GetJobID() string {
	return j.JobID
}
