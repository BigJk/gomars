package gomars

import (
	"io/ioutil"
	"math/rand"
)

// MARS is the simulator
type MARS struct {
	Coresize     int
	MaxProcess   int
	MaxCycles    int
	MaxLength    int
	MinDistance  int
	Core         Core
	Warrior      []Warrior
	WarriorCount int
}

// CreateMars creates a new MARS
func CreateMars(coresize, maxProcess, maxCycles, maxLength, minDistance int) MARS {
	return MARS{coresize, maxProcess, maxCycles, maxLength, minDistance, CreateCore(coresize), make([]Warrior, 0), 0}
}

// ParseWarrior ...
func (m *MARS) ParseWarrior(warrior string) Warrior {
	return ParseWarrior(m.Coresize, m.Core.SizePSpace, m.MaxCycles, m.MaxProcess, len(m.Warrior), m.MaxLength, m.MinDistance, warrior)
}

// ParseWarriorFromFile ...
func (m *MARS) ParseWarriorFromFile(path string) Warrior {
	w, _ := ioutil.ReadFile(path)
	return m.ParseWarrior(string(w))
}

// AddWarriorString adds a warrior from string
func (m *MARS) AddWarriorString(wstr string) {
	w := m.ParseWarrior(wstr)
	m.AddWarrior(w)
}

// AddWarrior adds a new warrior into the MARS
func (m *MARS) AddWarrior(w Warrior) {
	m.Warrior = append(m.Warrior, w)
	if w.HasPSpace() {
		m.Core.DisablePSpace = false
	}
}

// Clear cleans the MARS
func (m *MARS) Clear() {
	m.Warrior = make([]Warrior, 0)
	for i := 0; i < m.Coresize; i++ {
		m.Core.Memory[i] = Command{0, 0, 0, 0, 0, 0}
	}
	m.Core.Warriors = make([]*CoreWarrior, 0)
}

// Run executes the warrior in the core over multiple rounds
func (m *MARS) Run(rounds int) []int {
	m.WarriorCount = len(m.Warrior)
	wc := m.WarriorCount

	if !m.Core.DisablePSpace {
		m.Core.CreatePSpace(wc)
	}

	m.Core.Warriors = make([]*CoreWarrior, wc)
	w := make([]int, wc+1)
	for i := 0; i < rounds; i++ {
		for j := 0; j < m.Coresize; j++ {
			m.Core.Memory[j].Empty()
		}

		m.Core.PlaceWarrior(0, 0, m.MaxProcess, m.Warrior[0])
		m.Core.Warriors[0].ID = 0

		for j := 1; j < wc; j++ {
			m.Core.PlaceRandom(j, m.MinDistance, m.MaxProcess, m.Warrior[j])
			m.Core.Warriors[j].ID = j
		}

		winner := m.RunSingle(i)
		if !m.Core.DisablePSpace {
			for j := 0; j < wc; j++ {
				if winner != j {
					m.Core.PSpace[j][0] = 0
				} else {
					m.Core.PSpace[j][0] = m.Core.Alive
				}
			}
		}
		w[winner+1]++
	}
	return w
}

// RunSingle runs a singe round with a given offset to shift the starting warrior
func (m *MARS) RunSingle(offset int) int {
	wc := m.WarriorCount
	m.Core.Alive = wc

	for i := 0; i < m.MaxCycles; i++ {
		if m.Core.Alive <= 1 {
			break
		}
		for j := 0; j < wc; j++ {
			m.Core.Execute(m.Core.Warriors[(j+offset)%wc])
			if m.Core.Alive == 1 {
				break
			}
		}
	}

	if m.Core.Alive == wc {
		return -1
	}

	for j := 0; j < wc; j++ {
		if m.Core.Warriors[j].Alive() {
			return j
		}
	}

	return -1
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
