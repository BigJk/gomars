package gomars

// TaskQueue is a basic FIFO queue with a fixed size.
type TaskQueue struct {
	nodes []int
	size  int
	head  int
	tail  int
	count int
}

// NewTaskQueue returns a new queue with the given initial size.
func NewTaskQueue(size int) *TaskQueue {
	return &TaskQueue{
		nodes: make([]int, size),
		size:  size,
	}
}

// Push adds a node to the queue.
func (q *TaskQueue) Push(address int) {
	if q.count >= q.size {
		return
	}
	q.nodes[q.tail] = address
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *TaskQueue) Pop() int {
	if q.count == 0 {
		return -1
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}
