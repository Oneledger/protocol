/*

 */

package jobs

type Job interface {
	DoMyJob(ctx interface{})
	IsDone() bool

	GetType() string
	GetJobID() string
}

type Status int

const (
	New Status = iota
	InProgress
	Completed
	Failed

	Max_Retry_Count int = 3
)
