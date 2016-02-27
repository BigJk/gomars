package gomars

// MultiMARS ...
type MultiMARS struct {
	MARSs []MARS
}

// CreateMultiMars ...
func CreateMultiMars(coresize, maxProcess, cycles, maxLength, marsCount int) MultiMARS {
	m := MultiMARS{}
	for i := 0; i < marsCount; i++ {
		m.MARSs = append(m.MARSs, CreateMars(coresize, maxProcess, cycles, maxLength))
	}
	return m
}

// AddWarrior ...
func (m *MultiMARS) AddWarrior(commands []Command) {
	for i := 0; i < len(m.MARSs); i++ {
		m.MARSs[i].AddWarrior(commands)
	}
}

// Run ...
func (m *MultiMARS) Run(rounds int) map[int]int {
	rs := rounds / len(m.MARSs)
	o := rounds - (rs * len(m.MARSs))
	rsarr := make([]int, len(m.MARSs))

	for i := 0; i < len(m.MARSs); i++ {
		rsarr[i] += rs
	}

	i := 0
	for o > 0 {
		rsarr[i%len(m.MARSs)]++
		o--
		i++
	}

	r := make(map[int]int)
	rchan := make(chan map[int]int, len(m.MARSs))

	for i := 0; i < len(m.MARSs); i++ {
		go m.asyncRun(m.MARSs[i], rsarr[i], rchan)
	}

	for i := 0; i < len(m.MARSs); i++ {
		mr := <-rchan
		for key, val := range mr {
			r[key] += val
		}
	}

	return r
}

func (m *MultiMARS) asyncRun(mars MARS, rounds int, rchan chan map[int]int) {
	r := mars.Run(rounds)
	rchan <- r
}
