package gomars

import (
	"math/rand"
	"time"
)

// MARS is the simulator
type MARS struct {
	Coresize   int
	MaxProcess int
	Cycles     int
	MaxLength  int
	Core       Core
	warrior    [][]Command
}

// CreateMars creates a new MARS
func CreateMars(coresize, maxProcess, cycles, maxLength int) MARS {
	rand.Seed(time.Now().Unix())
	return MARS{coresize, maxProcess, cycles, maxLength, CreateCore(coresize), make([][]Command, 0)}
}

// AddWarrior adds a new warrior into the MARS
func (m *MARS) AddWarrior(commands []Command) {
	m.warrior = append(m.warrior, commands)
}

// PlaceWarrior places a warrior in the core
func (m *MARS) PlaceWarrior(position int, commands []Command) {
	m.Core.PlaceWarrior(position, commands)
	m.Core.Warriors[len(m.Core.Warriors)-1].Task = NewTaskQueue(m.MaxProcess)
	m.Core.Warriors[len(m.Core.Warriors)-1].Task.Push(position)
}

// ClearMars cleans the MARS
func (m *MARS) ClearMars() {
	m.warrior = make([][]Command, 0)
	for i := 0; i < m.Coresize; i++ {
		m.Core.Memory[i] = Command{0, 0, 0, 0, 0, 0}
	}
	m.Core.Warriors = make([]Warrior, 0)
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

		m.Core.Warriors = make([]Warrior, 0)

		m.PlaceWarrior(0, m.warrior[0])
		m.Core.Warriors[0].ID = 0

		for j := 1; j < len(m.warrior); j++ {
			pos := 0
			for !m.Core.IsEmpty(pos, len(m.warrior[j])) {
				pos = random(m.MaxLength, m.Coresize-len(m.warrior[j]))
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

	for i := 0; i < m.Cycles; i++ {
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
