/*

 */

package internal

type Job struct {
	Data        interface{}
	HandlerName string
}

func NewJob(name string, data interface{}) *Job {
	return &Job{
		Data:        data,
		HandlerName: name,
	}
}
