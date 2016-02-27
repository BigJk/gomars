package gomars

import "fmt"

var emptyCommand = Command{0, 0, 0, 0, 0, 0}

// Core is the virtual memory that the MARS uses
type Core struct {
	Warriors   []Warrior
	Memory     []Command
	PSpace     [][]int
	Size       int
	SizePSpace int
}

// CreateCore creates a new core
func CreateCore(size int) Core {
	c := Core{}
	c.Memory = make([]Command, size)
	c.Size = len(c.Memory)
	c.SizePSpace = c.Size / 16
	return c
}

// Print ...
func (c *Core) Print(s, e int) {
	fmt.Println("=====================")
	for i := 0; i < e; i++ {
		c.Memory[s+i].Print()
	}
	fmt.Println("=====================")
}

// Alive ...
func (c *Core) Alive() int {
	a := 0
	for i := 0; i < len(c.Warriors); i++ {
		if c.Warriors[i].Alive() {
			a++
		}
	}
	return a
}

// IsEmpty checks if a part of the core is empty
func (c *Core) IsEmpty(s, e int) bool {
	for i := 0; i < e; i++ {
		if !c.Memory[s+i].Equal(&emptyCommand) {
			return false
		}
	}
	return true
}

// PlaceWarrior places a warrior in the core
func (c *Core) PlaceWarrior(position int, commands []Command) {
	c.Warriors = append(c.Warriors, Warrior{})
	for i := 0; i < len(commands); i++ {
		c.Memory[position+i] = commands[i]
	}
}

// GetAddress ...
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

// GetAddresses ...
func (c *Core) GetAddresses(position int) (int, int, int, int) {
	a, apost := c.GetAddress(position, c.Memory[position].A, c.Memory[position].AddressingModeA)
	b, bpost := c.GetAddress(position, c.Memory[position].B, c.Memory[position].AddressingModeB)
	return a, b, apost, bpost
}

// NormalizeAddress ...
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

// NormalizePSpaceAddress ...
func (c *Core) NormalizePSpaceAddress(address int) int {
	if address >= 0 && address < c.SizePSpace {
		return address
	}
	r := address % c.SizePSpace
	if r < 0 {
		return c.SizePSpace + r
	}
	return r
}

// ExecuteCommand ...
func (c *Core) ExecuteCommand(w *Warrior, address int) {
	a, b, apost, bpost := c.GetAddresses(address)
	cmd := c.Memory[address]

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

}

// Execute ...
func (c *Core) Execute(w *Warrior) {
	if !w.Alive() {
		return
	}
	p := w.GetTask()
	c.ExecuteCommand(w, p)
}
