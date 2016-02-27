package gomars

// Task ...
type Task struct {
	Address int
}

// TaskQueue is a basic FIFO queue based on a circular list that resizes as needed.
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
	/*if q.head == q.tail && q.count > 0 {
		nodes := make([]int, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}*/
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
