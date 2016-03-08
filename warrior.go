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
	out := "PIN 0\nORG " + fmt.Sprint(w.EntryPoint)
	for i := 0; i < len(w.Code); i++ {
		out += "\n" + w.Code[i].ToString()
	}
	ioutil.WriteFile(path, []byte(out), 0777)
}
