/*

 */

package internal

import (
	"fmt"
)

type jobID string
type jobHandler func(data interface{})

type JobQueue struct {
	store    map[jobID]Job
	handlers map[string]jobHandler
}

func NewJobQueue() *JobQueue {
	return &JobQueue{
		store: make(map[jobID]Job),
	}
}

func (jq *JobQueue) DoJobs() {

	for id, job := range jq.store {
		fmt.Println(id)

		handler := jq.handlers[job.HandlerName]
		handler(job.Data)
	}
}
