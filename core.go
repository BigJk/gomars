package gomars

import "fmt"

var emptyCommand = Command{0, 0, 0, 0, 0, 0}

// Core is the virtual memory that the MARS uses
type Core struct {
	Warriors      []*CoreWarrior
	Memory        []Command
	PSpace        [][]int
	Size          int
	SizePSpace    int
	Alive         int
	DisablePSpace bool
}

// CreateCore creates a new core
func CreateCore(size int) Core {
	c := Core{}
	c.Memory = make([]Command, size)
	c.Size = len(c.Memory)
	c.SizePSpace = c.Size / 16
	c.DisablePSpace = true
	return c
}

// Print prints a section of the core
func (c *Core) Print(s, e int) {
	fmt.Println("=====================")
	for i := 0; i < e; i++ {
		c.Memory[s+i].Print()
	}
	fmt.Println("=====================")
}

// CreatePSpace creates the PSpace
func (c *Core) CreatePSpace(count int) {
	c.PSpace = make([][]int, count)
	for i := 0; i < count; i++ {
		c.PSpace[i] = make([]int, c.SizePSpace)
		c.PSpace[i][0] = -1
	}
}

// IsEmpty checks if a part of the core is empty
func (c *Core) IsEmpty(s, e int) bool {
	for i := 0; i < e; i++ {
		if !c.Memory[s+i].IsEmpty() {
			return false
		}
	}
	return true
}

// PlaceWarrior places a warrior in the core
func (c *Core) PlaceWarrior(num, position, maxProcess int, w Warrior) {
	c.Warriors[num] = &CoreWarrior{}
	for i := 0; i < len(w.Code); i++ {
		c.Memory[position+i] = w.Code[i]
	}
	c.Warriors[num].Task = NewTaskQueue(maxProcess)
	c.Warriors[num].Task.Push(position + w.EntryPoint)
}

// PlaceRandom ...
func (c *Core) PlaceRandom(num, minDistance, maxProcess int, w Warrior) {
	pos := random(minDistance, len(c.Memory)-len(w.Code))
	for !c.IsEmpty(pos, len(w.Code)) {
		pos = random(minDistance, len(c.Memory)-len(w.Code))
	}
	c.PlaceWarrior(num, pos, maxProcess, w)
}

// GetAddress returns the address and the postincrement address
func (c *Core) GetAddress(position, value int, addressingMode AddressingMode) (int, int) {
	vp := c.NormalizeAddress(position + value)

	switch addressingMode {
	case immediate:
		return position, 0
	case direct:
		return vp, 0
	case aIndirect:
		return c.NormalizeAddress(vp + c.Memory[vp].A), 0
	case bIndirect:
		return c.NormalizeAddress(vp + c.Memory[vp].B), 0
	case aIndirectPre:
		c.Memory[vp].A--
		return c.NormalizeAddress(vp + c.Memory[vp].A), 0
	case bIndirectPre:
		c.Memory[vp].B--
		return c.NormalizeAddress(vp + c.Memory[vp].B), 0
	case aIndirectPost:
		return c.NormalizeAddress(vp + c.Memory[vp].A), vp
	case bIndirectPost:
		return c.NormalizeAddress(vp + c.Memory[vp].B), vp
	}

	return 0, 0
}

// GetAddresses gets both addresses of the a and b field and also returns their addresses for postincrement
func (c *Core) GetAddresses(position int) (int, int, int, int) {
	cmd := &c.Memory[position]
	a, apost := c.GetAddress(position, cmd.A, cmd.AddressingModeA)
	b, bpost := c.GetAddress(position, cmd.B, cmd.AddressingModeB)
	return a, b, apost, bpost
}

// NormalizeAddress folds the address around the core
func (c *Core) NormalizeAddress(address int) int {
	if address >= 0 && address < c.Size {
		return address
	}
	r := address % c.Size
	if r < 0 {
		return c.Size + r
	}
	return r
}

// NormalizePSpaceAddress folds the address around the pspace
func (c *Core) NormalizePSpaceAddress(address int) int {
	if c.DisablePSpace {
		return 0
	}
	if address >= 0 && address < c.SizePSpace {
		return address
	}
	r := address % c.SizePSpace
	if r < 0 {
		return c.SizePSpace + r
	}
	return r
}

// ExecuteCommand executes a single command at the given address
func (c *Core) ExecuteCommand(w *CoreWarrior, address int) {
	a, b, apost, bpost := c.GetAddresses(address)
	cmd := &c.Memory[address]

	switch cmd.OpCode {
	case dat:
	case mov:
		c.mov(a, b, cmd.Modifier, address, w)
	case add:
		c.add(a, b, cmd.Modifier, address, w)
	case sub:
		c.sub(a, b, cmd.Modifier, address, w)
	case mul:
		c.mul(a, b, cmd.Modifier, address, w)
	case div:
		c.div(a, b, cmd.Modifier, address, w)
	case mod:
		c.mod(a, b, cmd.Modifier, address, w)
	case jmp:
		c.jmp(a, address, w)
	case jmz:
		c.jmz(a, b, address, w)
	case jmn:
		c.jmn(a, b, address, w)
	case djn:
		c.djn(a, b, cmd.Modifier, address, w)
	case spl:
		c.spl(a, address, w)
	case cmp, seq:
		c.seq(a, b, cmd.Modifier, address, w)
	case sne:
		c.sne(a, b, cmd.Modifier, address, w)
	case slt:
		c.slt(a, b, cmd.Modifier, address, w)
	case nop:
		w.QueueTask(c.NormalizeAddress(address + 1))
	case stp:
		c.stp(a, b, cmd.Modifier, address, w)
	case ldp:
		c.ldp(a, b, cmd.Modifier, address, w)
	}

	if cmd.AddressingModeA == aIndirectPost {
		c.Memory[apost].A++
	} else if cmd.AddressingModeA == bIndirectPost {
		c.Memory[apost].B++
	}

	if cmd.AddressingModeB == aIndirectPost {
		c.Memory[bpost].A++
	} else if cmd.AddressingModeB == bIndirectPost {
		c.Memory[bpost].B++
	}

	if !w.Alive() {
		c.Alive--
	}

}

// Execute executes the next warrior task
func (c *Core) Execute(w *CoreWarrior) {
	if !w.Alive() {
		return
	}
	p := w.GetTask()
	c.ExecuteCommand(w, p)
}
