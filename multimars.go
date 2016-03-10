package gomars

import (
	"io/ioutil"
	"sync"
)

// MultiMARS is the multi process simulator
type MultiMARS struct {
	Coresize     int
	MaxProcess   int
	MaxCycles    int
	MaxLength    int
	MinDistance  int
	Worker       int
	Warrior      []Warrior
	WarriorCount int
}

// CreateMultiMars creates a new MARS
func CreateMultiMars(coresize, maxProcess, maxCycles, maxLength, minDistance, worker int) MultiMARS {
	return MultiMARS{coresize, maxProcess, maxCycles, maxLength, minDistance, worker, make([]Warrior, 0), 0}
}

// ParseWarrior ...
func (m *MultiMARS) ParseWarrior(warrior string) Warrior {
	return ParseWarrior(m.Coresize, m.Coresize/16, m.MaxCycles, m.MaxProcess, len(m.Warrior), m.MaxLength, m.MinDistance, warrior)
}

// ParseWarriorFromFile ...
func (m *MultiMARS) ParseWarriorFromFile(path string) Warrior {
	w, _ := ioutil.ReadFile(path)
	return m.ParseWarrior(string(w))
}

// AddWarriorString adds a warrior from string
func (m *MultiMARS) AddWarriorString(wstr string) {
	w := m.ParseWarrior(wstr)
	m.AddWarrior(w)
}

// AddWarrior adds a new warrior into the MARS
func (m *MultiMARS) AddWarrior(w Warrior) {
	m.Warrior = append(m.Warrior, w)
}

// Run ...
func (m *MultiMARS) Run(rounds int) []int {
	var wg sync.WaitGroup
	wg.Add(m.Worker)

	m.WarriorCount = len(m.Warrior)
	wc := m.WarriorCount
	w := make([]int, wc+1)

	r := make(chan int, rounds)
	in := make(chan int, rounds+m.Worker)

	for i := 0; i < rounds; i++ {
		in <- i
	}

	for i := 0; i < m.Worker; i++ {
		in <- -1
	}

	cores := make([]Core, m.Worker)
	for i := 0; i < m.Worker; i++ {
		cores[i] = CreateCore(m.Coresize)
		cores[i].Warriors = make([]*CoreWarrior, wc)
		go m.CoreWorker(&cores[i], in, r, &wg)
	}

	wg.Wait()

	for i := 0; i < rounds; i++ {
		w[<-r+1]++
	}

	return w
}

// CoreWorker ...
func (m *MultiMARS) CoreWorker(c *Core, in chan int, r chan int, wg *sync.WaitGroup) {
	offset := 0
	wc := m.WarriorCount

	for {

		offset = <-in
		if offset == -1 {
			wg.Done()
			return
		}

		c.Alive = wc

		for j := 0; j < m.Coresize; j++ {
			c.Memory[j].Empty()
		}

		c.PlaceWarrior(0, 0, m.MaxProcess, m.Warrior[0])
		c.Warriors[0].ID = 0

		for j := 1; j < wc; j++ {
			c.PlaceRandom(j, m.MinDistance, m.MaxProcess, m.Warrior[j])
			c.Warriors[j].ID = j
		}

		for i := 0; i < m.MaxCycles; i++ {
			if c.Alive <= 1 {
				break
			}
			for j := 0; j < wc; j++ {
				c.Execute(c.Warriors[(j+offset)%wc])
				if c.Alive == 1 {
					break
				}
			}
		}

		if c.Alive == wc {
			r <- -1
		} else {
			found := false
			for j := 0; j < wc; j++ {
				if c.Warriors[j].Alive() {
					r <- j
					found = true
					break
				}
			}
			if !found {
				r <- -1
			}
		}

	}
}
