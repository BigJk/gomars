package gomars

// Warrior represents a warrior in the core
type Warrior struct {
	ID   int
	Task *TaskQueue
}

// GetTask ...
func (w *Warrior) GetTask() int {
	return w.Task.Pop()
}

// QueueTask ...
func (w *Warrior) QueueTask(address int) {
	w.Task.Push(address)
}

// Alive ...
func (w *Warrior) Alive() bool {
	return w.Task.count != 0
}
