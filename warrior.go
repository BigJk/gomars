package gomars

// Warrior represents a warrior in the core
type Warrior struct {
	ID   int
	Task *TaskQueue
}

// GetTask gets the next task
func (w *Warrior) GetTask() int {
	return w.Task.Pop()
}

// QueueTask queues the next task
func (w *Warrior) QueueTask(address int) {
	w.Task.Push(address)
}

// Alive checks if the warrior is still alive
func (w *Warrior) Alive() bool {
	return w.Task.count != 0
}
