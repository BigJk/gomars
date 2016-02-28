package gomars

import (
	"regexp"
	"strconv"
	"strings"
)

var lineRegex = regexp.MustCompile("([a-z]+).([a-z]+)[ ]*(.)[ ]*([-0-9]+)[ ]*,[ ]*(.)[ ]*([-0-9]+)")

// ParseWarrior ...
func ParseWarrior(s string) Warrior {
	var cmds []Command
	org := 0
	sp := strings.Split(s, "\n")
	for i := 0; i < len(sp); i++ {
		if sp[i] == "" || strings.HasPrefix(sp[i], "PIN") {
			continue
		}
		if strings.HasPrefix(sp[i], "ORG") {
			org, _ = strconv.Atoi(sp[i][4:])
			continue
		}
		line, valid := ParseLine(sp[i])
		if valid {
			cmds = append(cmds, line)
		}
	}
	return Warrior{org, cmds}
}

// ParseLine ...
func ParseLine(s string) (Command, bool) {
	m := lineRegex.FindAllStringSubmatch(strings.ToLower(s), 1)[0]
	if len(m) != 7 {
		return emptyCommand, false
	}
	a, _ := strconv.Atoi(m[4])
	b, _ := strconv.Atoi(m[6])
	return Command{ParseOpCode(m[1]), ParseModifier(m[2]), ParseAddressingMode(m[3]), a, ParseAddressingMode(m[5]), b}, true
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
