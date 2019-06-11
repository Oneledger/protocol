package identity

import (
	"container/heap"

	"github.com/Oneledger/protocol/utils"
)

type ValidatorQueue struct {
	utils.PriorityQueue
}

func (vq *ValidatorQueue) Pop() *utils.Queued {
	if vq.Len() < 1 {
		return nil
	}
	return heap.Pop(&vq.PriorityQueue).(*utils.Queued)
}

func (vq *ValidatorQueue) Push(queued *utils.Queued) {
	heap.Push(&vq.PriorityQueue, queued)
}

func (vq *ValidatorQueue) Init() {
	heap.Init(&vq.PriorityQueue)
}

func (vq *ValidatorQueue) Len() int {
	return vq.PriorityQueue.Len()
}

func (vq *ValidatorQueue) append(queued *utils.Queued) {
	vq.PriorityQueue = append(vq.PriorityQueue, queued)
}

func (vq *ValidatorQueue) update(queued *utils.Queued, value []byte, priority int64) {
	vq.PriorityQueue.Update(queued, value, priority)
}
