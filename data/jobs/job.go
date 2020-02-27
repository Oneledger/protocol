/*

 */

package jobs

type Job interface {
	DoMyJob(ctx interface{})
	IsDone() bool
	IsFailed() bool
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

type RedeemStatus int8

const (
	NewRedeem              = -1
	Ongoing   RedeemStatus = iota
	Success
	Expired
)
