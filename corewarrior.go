package gomars

// CoreWarrior represents a warrior in the core
type CoreWarrior struct {
	ID   int
	Task *TaskQueue
}

// GetTask gets the next task
func (w *CoreWarrior) GetTask() int {
	return w.Task.Pop()
}

// QueueTask queues the next task
func (w *CoreWarrior) QueueTask(address int) {
	w.Task.Push(address)
}

// Alive checks if the warrior is still alive
func (w *CoreWarrior) Alive() bool {
	return w.Task.Count != 0
}
