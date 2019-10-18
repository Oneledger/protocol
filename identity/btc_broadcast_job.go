/*

 */

package identity

type JobBTCBroadcast struct {
	Type string

	TrackerName string

	JobID string

	Done bool
}

func (j *JobBTCBroadcast) DoMyJob(ctx *JobsContext) {
	panic("implement me")
}

func (j *JobBTCBroadcast) IsMyJobDone(ctx *JobsContext) bool {
	panic("implement me")
}

func (j *JobBTCBroadcast) IsSufficient() bool {
	panic("implement me")
}

func (j *JobBTCBroadcast) DoFinalize() {
	panic("implement me")
}

/*
	simple getters
*/
func (j *JobBTCBroadcast) GetType() string {
	return JobTypeBTCBroadcast
}

func (j *JobBTCBroadcast) GetJobID() string {
	return j.JobID
}

func (j *JobBTCBroadcast) IsDone() bool {
	return j.Done
}
