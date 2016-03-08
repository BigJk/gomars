package gomars

import "fmt"

// OpCode ...
type OpCode byte

const (
	dat OpCode = iota // data (kills the process)
	mov               // move (copies data from one address to another)
	add               // add (adds one number to another)
	sub               // subtract (subtracts one number from another)
	mul               // multiply (multiplies one number with another)
	div               // divide (divides one number with another)
	mod               // modulus (divides one number with another and gives the remainder)
	jmp               // jump (continues execution from another address)
	jmz               // jump if zero (tests a number and jumps to an address if it's 0)
	jmn               // jump if not zero (tests a number and jumps if it isn't 0)
	djn               // decrement and jump if not zero (decrements a number by one, and jumps unless the result is 0)
	spl               // split (starts a second process at another address)
	cmp               // compare (same as SEQ)
	seq               // skip if equal (compares two instructions, and skips the next instruction if they are equal)
	sne               // skip if not equal (compares two instructions, and skips the next instruction if they aren't equal)
	slt               // skip if lower than (compares two values, and skips the next instruction if the first is lower than the second)
	ldp               // load from p-space (loads a number from private storage space)
	stp               // save to p-space (saves a number to private storage space)
	nop               // no operation (does nothing)
)

// Modifier ...
type Modifier byte

const (
	f  Modifier = iota // moves both fields of the source into the same fields in the destination
	a                  // moves the A-field of the source into the A-field of the destination
	b                  // moves the B-field of the source into the B-field of the destination
	ab                 // moves the A-field of the source into the B-field of the destination
	ba                 // moves the B-field of the source into the A-field of the destination
	x                  // moves both fields of the source into the opposite fields in the destination
	i                  // moves the whole source instruction into the destination
)

// AddressingMode ...
type AddressingMode byte

const (
	immediate     AddressingMode = iota // immediate
	direct                              // direct (the $ may be omitted)
	aIndirect                           // A-field indirect
	bIndirect                           // B-field indirect
	aIndirectPre                        // A-field indirect with predecrement
	bIndirectPre                        // B-field indirect with predecrement
	aIndirectPost                       // A-field indirect with postincrement
	bIndirectPost                       // B-field indirect with postincrement
)

// Command is a single line of redcode
type Command struct {
	OpCode          OpCode
	Modifier        Modifier
	AddressingModeA AddressingMode
	A               int
	AddressingModeB AddressingMode
	B               int
}

// Print prints a command
func (c *Command) Print() {
	fmt.Println(c.OpCode.String(), c.Modifier.String(), c.AddressingModeA.String(), c.A, c.AddressingModeB.String(), c.B)
}

// ToString ...
func (c *Command) ToString() string {
	return c.OpCode.String() + "." + c.Modifier.String() + " " + c.AddressingModeA.String() + fmt.Sprint(c.A) + ", " + c.AddressingModeB.String() + fmt.Sprint(c.B)
}

// Equal checks if two commands are equal
func (c *Command) Equal(o *Command) bool {
	return c.OpCode == o.OpCode && c.Modifier == o.Modifier && c.AddressingModeA == o.AddressingModeA && c.A == o.A && c.AddressingModeB == o.AddressingModeB && c.B == o.B
}

// IsEmpty ...
func (c *Command) IsEmpty() bool {
	return (c.A == 0 && c.B == 0 && c.AddressingModeA == 0 && c.AddressingModeB == 0 && c.Modifier == 0 && c.OpCode == 0)
}

// Clone clones the command
func (c *Command) Clone() Command {
	return Command{c.OpCode, c.Modifier, c.AddressingModeA, c.A, c.AddressingModeB, c.B}
}

// Empty emptys the command
func (c *Command) Empty() {
	c.OpCode = 0
	c.Modifier = 0
	c.A = 0
	c.B = 0
	c.AddressingModeA = 0
	c.AddressingModeB = 0
}
