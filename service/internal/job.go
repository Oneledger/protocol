/*

 */

package internal

type Job interface {
	DoProcess()

	IsDone()
	IsRequeue()
	IsFailed()

	DoSuccess()
	DoFailure()
	DoRequeue()
}

type Job struct {
	Data        interface{}
	HandlerName string
	Result      interface{}
}

func NewJob(name string, data interface{}) *Job {
	return &Job{
		Data:        data,
		HandlerName: name,
	}
}
