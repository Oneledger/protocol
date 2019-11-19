/*

 */

package jobs

type Job interface {
	DoMyJob(ctx interface{})

	GetType() string
	GetJobID() string
}
