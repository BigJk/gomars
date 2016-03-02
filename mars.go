package gomars

import (
	"math/rand"
	"time"
)

// MARS is the simulator
type MARS struct {
	Coresize     int
	MaxProcess   int
	MaxCycles    int
	MaxLength    int
	MinDistance  int
	Core         Core
	warrior      []Warrior
	warriorCount int
}

// CreateMars creates a new MARS
func CreateMars(coresize, maxProcess, maxCycles, maxLength, minDistance int) MARS {
	rand.Seed(time.Now().Unix())
	return MARS{coresize, maxProcess, maxCycles, maxLength, minDistance, CreateCore(coresize), make([]Warrior, 0), 0}
}

// AddWarriorString adds a warrior from string
func (m *MARS) AddWarriorString(wstr string) {
	w := m.ParseWarrior(wstr)
	m.AddWarrior(w)
}

// AddWarrior adds a new warrior into the MARS
func (m *MARS) AddWarrior(w Warrior) {
	m.warrior = append(m.warrior, w)
	if w.HasPSpace() {
		m.Core.DisablePSpace = false
	}
}

// PlaceWarrior places a warrior in the core
func (m *MARS) PlaceWarrior(num, position int, w Warrior) {
	m.Core.PlaceWarrior(num, position, w.Code)
	m.Core.Warriors[num].Task = NewTaskQueue(m.MaxProcess)
	m.Core.Warriors[num].Task.Push(position + w.EntryPoint)
}

// PlaceRandom places a warrior randomly in the core
func (m *MARS) PlaceRandom(num int, w Warrior) {
	pos := random(m.MinDistance, m.Coresize-len(w.Code))
	for !m.Core.IsEmpty(pos, len(w.Code)) {
		pos = random(m.MinDistance, m.Coresize-len(w.Code))
	}
	m.PlaceWarrior(num, pos, w)
}

// Clear cleans the MARS
func (m *MARS) Clear() {
	m.warrior = make([]Warrior, 0)
	for i := 0; i < m.Coresize; i++ {
		m.Core.Memory[i] = Command{0, 0, 0, 0, 0, 0}
	}
	m.Core.Warriors = make([]CoreWarrior, 0)
}

// Run executes the warrior in the core over multiple rounds
func (m *MARS) Run(rounds int) []int {
	m.warriorCount = len(m.warrior)
	wc := m.warriorCount

	if !m.Core.DisablePSpace {
		m.Core.CreatePSpace(wc)
	}

	m.Core.Warriors = make([]CoreWarrior, wc)
	w := make([]int, wc+1)
	for i := 0; i < rounds; i++ {
		for j := 0; j < m.Coresize; j++ {
			m.Core.Memory[j].Empty()
		}

		m.PlaceWarrior(0, 0, m.warrior[0])
		m.Core.Warriors[0].ID = 0

		for j := 1; j < wc; j++ {
			m.PlaceRandom(j, m.warrior[j])
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
	wc := m.warriorCount
	m.Core.Alive = wc

	for i := 0; i < m.MaxCycles; i++ {
		if m.Core.Alive <= 1 {
			break
		}
		for j := 0; j < wc; j++ {
			m.Core.Execute(&m.Core.Warriors[(j+offset)%wc])
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
