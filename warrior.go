package gomars

import (
	"fmt"
	"io/ioutil"
)

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

// SaveLoadFile ...
func (w *Warrior) SaveLoadFile(path string) {
	ioutil.WriteFile(path, []byte(w.LoadFile()), 0777)
}

// LoadFile ...
func (w *Warrior) LoadFile() string {
	out := "ORG " + fmt.Sprint(w.EntryPoint) + "\nPIN 0"
	for i := 0; i < len(w.Code); i++ {
		out += "\n" + w.Code[i].ToString()
	}
	return out
}
