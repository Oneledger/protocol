package event

import (
	"strconv"

	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"
)

var _ jobs.Job = &JobGovFinalizeProposal{}

type JobGovFinalizeProposal struct {
	ProposalID governance.ProposalID
	JobID      string
	RetryCount int
	Status     jobs.Status
}

func NewGovFinalizeProposalJob(proposalID governance.ProposalID, status governance.ProposalStatus) *JobGovFinalizeProposal {
	return &JobGovFinalizeProposal{
		ProposalID: proposalID,
		JobID:      proposalID.String() + storage.DB_PREFIX + strconv.Itoa(int(status)),
		RetryCount: 0,
		Status:     0,
	}
}

func (j *JobGovFinalizeProposal) DoMyJob(ctx interface{}) {
	govCtx, _ := ctx.(*JobsContext)
	BroadcastGovFinalizeVotesTx(govCtx, j.ProposalID, j.JobID)
	j.Status = jobs.Completed
}

func (j *JobGovFinalizeProposal) IsDone() bool {
	return j.Status == jobs.Completed
}

func (j *JobGovFinalizeProposal) IsFailed() bool {
	return j.Status == jobs.Failed
}

func (j *JobGovFinalizeProposal) GetType() string {
	return JobTypeGOVFinalizeProposal
}

func (j *JobGovFinalizeProposal) GetJobID() string {
	return j.JobID
}
