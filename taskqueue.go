package gomars

// TaskQueue is a basic FIFO queue with a fixed size.
type TaskQueue struct {
	Nodes []int
	Size  int
	Head  int
	Tail  int
	Count int
}

// NewTaskQueue returns a new queue with the given initial size.
func NewTaskQueue(size int) *TaskQueue {
	return &TaskQueue{
		Nodes: make([]int, size),
		Size:  size,
	}
}

// Push adds a node to the queue.
func (q *TaskQueue) Push(address int) {
	if q.Count >= q.Size {
		return
	}
	q.Nodes[q.Tail] = address
	q.Tail = (q.Tail + 1) % q.Size
	q.Count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *TaskQueue) Pop() int {
	if q.Count == 0 {
		return -1
	}
	node := q.Nodes[q.Head]
	q.Head = (q.Head + 1) % q.Size
	q.Count--
	return node
}
