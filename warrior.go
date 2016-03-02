package gomars

// Warrior ...
type Warrior struct {
	EntryPoint int
	Code       []Command
}

// HasPSpace ...
func (w *Warrior) HasPSpace() bool {
	for i := 0; i < len(w.Code); i++ {
		if w.Code[i].OpCode == stp || w.Code[i].OpCode == ldp {
			return true
		}
	}
	return false
}
