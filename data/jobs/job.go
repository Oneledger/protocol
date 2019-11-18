/*

 */

package jobs

type Job interface {
	DoMyJob(ctx interface{})
	IsMyJobDone(ctx interface{}) bool

	IsSufficient(ctx interface{}) bool
	DoFinalize()

	GetType() string
	GetJobID() string
	IsDone() bool
}


type Status int
const (
	New  Status = iota
	InProgress
	Completed
	Failed




	Max_Retry_Count int = 3
)

