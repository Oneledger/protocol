package identity

import "container/heap"

// An Queued is a priority queue we hold for validators
type Queued struct {
	value    []byte // The value of the item; arbitrary.
	priority int64  // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A ValidatorQueue implements heap.Interface and holds Items.
type ValidatorQueue []*Queued

func (vq ValidatorQueue) Len() int { return len(vq) }

func (vq ValidatorQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return vq[i].priority > vq[j].priority
}

func (vq ValidatorQueue) Swap(i, j int) {
	vq[i], vq[j] = vq[j], vq[i]
	vq[i].index = i
	vq[j].index = j
}

func (vq *ValidatorQueue) Push(x interface{}) {
	n := len(*vq)
	item := x.(*Queued)
	item.index = n
	*vq = append(*vq, item)
}

func (vq *ValidatorQueue) Pop() interface{} {
	old := *vq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*vq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Queued in the queue.
func (vq *ValidatorQueue) update(item *Queued, value []byte, priority int64) {
	item.value = value
	item.priority = priority
	heap.Fix(vq, item.index)
}
