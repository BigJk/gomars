package gomars

import (
	"math/rand"
	"time"
)

// Warrior ...
type Warrior struct {
	EntryPoint int
	Code       []Command
}

// MARS is the simulator
type MARS struct {
	Coresize    int
	MaxProcess  int
	MaxCycles   int
	MaxLength   int
	MinDistance int
	Core        Core
	warrior     []Warrior
}

// CreateMars creates a new MARS
func CreateMars(coresize, maxProcess, maxCycles, maxLength, minDistance int) MARS {
	rand.Seed(time.Now().Unix())
	return MARS{coresize, maxProcess, maxCycles, maxLength, minDistance, CreateCore(coresize), make([]Warrior, 0)}
}

// AddWarriorString adds a warrior from string
func (m *MARS) AddWarriorString(wstr string) {
	w := ParseWarrior(m.ParseToLoadFile(wstr))
	m.AddWarrior(w)
}

// AddWarrior adds a new warrior into the MARS
func (m *MARS) AddWarrior(w Warrior) {
	m.warrior = append(m.warrior, w)
}

// PlaceWarrior places a warrior in the core
func (m *MARS) PlaceWarrior(position int, w Warrior) {
	m.Core.PlaceWarrior(position, w.Code)
	m.Core.Warriors[len(m.Core.Warriors)-1].Task = NewTaskQueue(m.MaxProcess)
	m.Core.Warriors[len(m.Core.Warriors)-1].Task.Push(position + w.EntryPoint)
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
func (m *MARS) Run(rounds int) map[int]int {
	m.Core.PSpace = make([][]int, len(m.warrior))
	for i := 0; i < len(m.warrior); i++ {
		m.Core.PSpace[i] = make([]int, m.Core.SizePSpace)
		m.Core.PSpace[i][0] = -1
	}

	w := make(map[int]int)
	for i := 0; i < rounds; i++ {
		for j := 0; j < m.Coresize; j++ {
			m.Core.Memory[j] = Command{0, 0, 0, 0, 0, 0}
		}

		m.Core.Warriors = make([]CoreWarrior, 0)

		m.PlaceWarrior(0, m.warrior[0])
		m.Core.Warriors[0].ID = 0

		for j := 1; j < len(m.warrior); j++ {
			pos := 0
			for !m.Core.IsEmpty(pos, len(m.warrior[j].Code)) {
				pos = random(m.MinDistance, m.Coresize-len(m.warrior[j].Code))
			}
			m.PlaceWarrior(pos, m.warrior[j])
			m.Core.Warriors[j].ID = j
		}

		winner := m.RunSingle(i)
		for i := 0; i < len(m.warrior); i++ {
			if winner != i {
				m.Core.PSpace[i][0] = 0
			} else {
				m.Core.PSpace[i][0] = m.Core.Alive()
			}
		}
		w[winner]++
	}
	return w
}

// RunSingle runs a singe round with a given offset to shift the starting warrior
func (m *MARS) RunSingle(offset int) int {
	wl := len(m.Core.Warriors)

	for i := 0; i < m.MaxCycles; i++ {
		if m.Core.Alive() <= 1 {
			break
		}
		for j := 0; j < wl; j++ {
			if m.Core.Warriors[(j+offset)%wl].Alive() {
				m.Core.Execute(&m.Core.Warriors[(j+offset)%wl])
			}
			if m.Core.Alive() == 1 {
				break
			}
		}
	}

	if m.Core.Alive() == len(m.warrior) {
		return -1
	}

	for j := 0; j < len(m.Core.Warriors); j++ {
		if m.Core.Warriors[j].Alive() {
			return j
		}
	}

	return -1
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
