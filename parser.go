package gomars

import (
	"regexp"
	"strconv"
	"strings"
)

var lineRegex = regexp.MustCompile("([a-z]+).([a-z]+)[ ]*(.)[ ]*([-0-9]+),[ ]*(.)[ ]*([-0-9]+)")

// ParseWarrior ...
func ParseWarrior(s string) []Command {
	var cmds []Command
	sp := strings.Split(s, "\n")
	for i := 0; i < len(sp); i++ {
		cmds = append(cmds, ParseLine(sp[i]))
	}
	return cmds
}

// ParseLine ...
func ParseLine(s string) Command {
	m := lineRegex.FindAllStringSubmatch(strings.ToLower(s), 1)[0]
	a, _ := strconv.Atoi(m[4])
	b, _ := strconv.Atoi(m[6])
	return Command{ParseOpCode(m[1]), ParseModifier(m[2]), ParseAddressingMode(m[3]), a, ParseAddressingMode(m[5]), b}
}

// ParseOpCode ...
func ParseOpCode(s string) OpCode {
	for i := 0; i < int(nop+1); i++ {
		if OpCode(i).String() == s {
			return OpCode(i)
		}
	}
	return 0
}

// ParseModifier ...
func ParseModifier(s string) Modifier {
	for i := 0; i < int(i+1); i++ {
		if Modifier(i).String() == s {
			return Modifier(i)
		}
	}
	return 0
}

// ParseAddressingMode ...
func ParseAddressingMode(s string) AddressingMode {
	switch s {
	case "#":
		return immediate
	case "$":
		return direct
	case "*":
		return aIndirect
	case "@":
		return bIndirect
	case "{":
		return aIndirectPre
	case "<":
		return bIndirectPre
	case "}":
		return aIndirectPost
	case ">":
		return bIndirectPost
	}
	return 0
}
