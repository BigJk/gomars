package gomars

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/sanderhahn/goexpr/eval"
)

type stringSorter []string

func (s stringSorter) Len() int {
	return len(s)
}
func (s stringSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s stringSorter) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

var comma = regexp.MustCompile("[ ]*,[ ]*")
var comments = regexp.MustCompile(";.*\n")
var emptySpace = regexp.MustCompile("[ 	]+")
var isFormular = regexp.MustCompile("[0-9][+\\-*\\/%]+[0-9]")
var field = regexp.MustCompile("^([*#{}$<>@])?(.*)")
var addressingModes = regexp.MustCompile("([a-z]{3}) ([*#{}$<>@]).*, ([*#{}$<>@])")
var isNoNumber = regexp.MustCompile("[a-zA-Z]+")
var singleLine = regexp.MustCompile("([a-z]+)\\.?([a-z]+)?[ ]*(.)[ ]*([-0-9]+)[ ]*,[ ]*(.)[ ]*([-0-9]+)")
var surround = regexp.MustCompile("\\((-?[0-9]+)\\)")

var opcodes = []string{"dat", "mov", "add", "sub", "mul", "div", "mod", "jmp", "jmz", "jmn", "djn", "spl", "seq", "sne", "cmp", "slt", "nop", "ldp", "stp"}
var modifier = []string{"*", "#", "{", "}", "$", ">", "<", "@"}
var env = eval.NewEnvironment()

// ParseWarrior ...
func (m *MARS) ParseWarrior(warrior string) Warrior {
	predefined := make(map[string]int)
	predefined["CORESIZE"] = m.Coresize
	predefined["PSPACESIZE"] = len(m.Core.PSpace)
	predefined["MAXCYCLES"] = m.MaxCycles
	predefined["MAXPROCESSES"] = m.MaxProcess
	predefined["WARRIORS"] = len(m.Warrior) // TODO: FIX!
	predefined["MAXLENGTH"] = m.MaxLength
	predefined["MINDISTANCE"] = m.MinDistance

	// CONVERT TO STRING AND NORMALIZE CODE
	warrior = normalizeCode(warrior)
	warrior = inserPredefined(warrior, predefined)

	// SPLIT WARRIOR
	warriorLines := strings.Split(warrior, "\n")
	var compiledWarriorLines []string

	curline := 0
	labels := make(map[string]int)
	equs := make(map[string]string)
	org := ""
	end := ""
	inFor := false

	for i := 0; i < len(warriorLines); i++ {
		if inFor {
			if strings.ToLower(warriorLines[i]) == "rof" {
				inFor = false
			}
			continue
		}

		if warriorLines[i] == "" {
			continue
		}
		split := strings.Split(warriorLines[i], " ")

		if len(split) == 1 {
			labels[split[0]] = curline
			continue
		}

		switch strings.ToLower(split[0]) {
		case "org":
			if len(split) == 2 {
				org = split[1]
			}
			continue
		case "end":
			if len(split) == 2 {
				end = split[1]
			}
			continue
		}

		if !isOpcode(split[0]) {
			if isFor(split[0]) {
				c, _ := env.Eval(replaceCurline(split[1], curline))
				inline := 0
				for {
					if strings.ToLower(warriorLines[i+inline]) == "rof" {
						break
					}
					inline++
				}
				for j := 0; j < int(c); j++ {
					for k := 0; k < inline-1; k++ {
						compiledWarriorLines = append(compiledWarriorLines, warriorLines[i+1+k])
						curline++
					}
				}
				inFor = true
			} else if split[1] == "equ" {
				equs[split[0]] = split[2]
			} else {
				if isOpcode(split[1]) {
					labels[split[0]] = curline
					compiledWarriorLines = append(compiledWarriorLines, warriorLines[i][len(split[0])+1:])
					curline++
				}
			}
		} else {
			compiledWarriorLines = append(compiledWarriorLines, warriorLines[i])
			curline++
		}

	}

	compiledWarrior := Warrior{}

	var labelKeys []string
	var equKeys []string

	for key := range equs {
		equKeys = append(equKeys, key)
	}

	for key := range labels {
		labelKeys = append(labelKeys, key)
	}

	sort.Sort(stringSorter(equKeys))
	sort.Sort(stringSorter(labelKeys))

	compiledWarrior.Code = make([]Command, len(compiledWarriorLines))

	for i := 0; i < len(compiledWarriorLines); i++ {

		split := strings.Split(compiledWarriorLines[i], " ")

		for j := 0; j < len(equKeys); j++ {
			split[1] = strings.Replace(split[1], equKeys[j], equs[equKeys[j]], -1)
			if len(split) >= 3 {
				split[2] = strings.Replace(split[2], equKeys[j], equs[equKeys[j]], -1)
			}
		}

		for j := 0; j < len(labelKeys); j++ {
			split[1] = strings.Replace(split[1], labelKeys[j], fmt.Sprint(labels[labelKeys[j]]-i), -1)
			if len(split) >= 3 {
				split[2] = strings.Replace(split[2], labelKeys[j], fmt.Sprint(labels[labelKeys[j]]-i), -1)
			}
		}

		var newLine bytes.Buffer

		if len(split) >= 3 {
			newLine.WriteString(split[0])
			newLine.WriteString(" ")
			newLine.WriteString(compileField(split[1], curline))
			newLine.WriteString(", ")
			newLine.WriteString(compileField(split[2], curline))
		} else {
			newLine.WriteString(split[0])
			newLine.WriteString(" ")
			newLine.WriteString(compileField(split[1], curline))
			newLine.WriteString(", $0")
		}

		compiledWarrior.Code[i], _ = ParseLine(newLine.String())

	}

	if org == "" && end != "" {
		org = end
	}

	if org != "" {
		if isNoNumber.MatchString(org) {
			for j := 0; j < len(labelKeys); j++ {
				org = strings.Replace(org, labelKeys[j], fmt.Sprint(labels[labelKeys[j]]), -1)
			}
			v, _ := env.Eval(org)
			org = fmt.Sprint(v)
		}
	} else {
		org = "0"
	}

	compiledWarrior.EntryPoint, _ = strconv.Atoi(org)

	return compiledWarrior
}

func inserPredefined(s string, predefined map[string]int) string {
	out := s
	for key, val := range predefined {
		out = strings.Replace(out, key, strconv.Itoa(val), -1)
	}
	return out
}

func normalizeCode(s string) string {
	// TODO: FOR 0 ROF COMMENTS STRIP
	out := emptySpace.ReplaceAllString(s, " ")
	out = comments.ReplaceAllString(out, "\n")
	out = strings.Replace(out, "\r", "", -1)
	out = strings.Replace(out, "\n ", "\n", -1)
	out = strings.Replace(out, ":", "", -1)
	out = comma.ReplaceAllString(out, " ")
	for i := 0; i < len(modifier); i++ {
		out = strings.Replace(out, modifier[i]+" ", modifier[i], -1)
	}
	return stripEmptyLines(out)
}

func stripEmptyLines(s string) string {
	split := strings.Split(s, "\n")
	out := ""
	for i := 0; i < len(split); i++ {
		if split[i] != " " && split[i] != "" {
			out += split[i] + "\n"
		}
	}
	return out[:len(out)-1]
}

func isOpcode(s string) bool {
	if len(s) < 3 {
		return false
	}
	if len(s) > 3 {
		if s[3:4] != "." {
			return false
		}
	}
	for i := 0; i < len(opcodes); i++ {
		if strings.ToLower(s)[:3] == opcodes[i] {
			return true
		}
	}
	return false
}

func fixForumlar(s string) string {
	out := strings.Replace(s, "(-", "(0-", -1)
	out = strings.Replace(s, "+-", "-", -1)
	m := surround.FindAllStringSubmatch(out, -1)
	for i := 0; i < len(m); i++ {
		if m[i][1][0] != '-' {
			out = strings.Replace(out, m[i][0], m[i][1], -1)
		} else {
			out = strings.Replace(out, m[i][0], "(0"+m[i][1]+")", -1)
		}
	}
	if out[0:1] == "-" {
		out = "0" + out
	}
	return out
}

func replaceCurline(s string, i int) string {
	return strings.Replace(s, "CURLINE", strconv.Itoa(i), -1)
}

func compileField(s string, i int) string {
	f := field.FindStringSubmatch(strings.Replace(strings.ToLower(s), "curline", fmt.Sprint(i), -1))
	//fmt.Println(s, f)
	if len(f) == 0 {
		if len(isFormular.FindAllString(s, -1)) == 0 {
			return "$" + s
		}
		v, _ := env.Eval(fixForumlar(s))
		return "$" + fmt.Sprint(math.Ceil(v))
	}
	if len(isFormular.FindAllString(s, -1)) == 0 {
		return s
	}
	v, _ := env.Eval(fixForumlar(f[2]))
	return f[1] + fmt.Sprint(math.Ceil(v))
}

func isFor(s string) bool {
	return strings.ToLower(s) == "for"
}

func opcodeStandard(o string, aAddr string, bAddr string) Modifier {
	if o == "dat" {
		return f
	}

	if (o == "mov" || o == "cmp") && aAddr == "#" || (o == "add" || o == "sub" || o == "mul" || o == "div" || o == "mod" || o == "slt") && aAddr == "#" {
		return ab
	} else if (o == "mov" || o == "cmp") && bAddr == "#" || (o == "add" || o == "sub" || o == "mul" || o == "div" || o == "mod") && aAddr == "#" {
		return b
	} else if o == "mov" || o == "cmp" {
		return i
	} else if o == "add" || o == "sub" || o == "mul" || o == "div" || o == "mod" {
		return f
	}

	return b
}

// ParseLine ...
func ParseLine(s string) (Command, bool) {
	m := singleLine.FindStringSubmatch(strings.ToLower(s))
	if len(m) != 7 && len(m) != 6 {
		return emptyCommand, false
	}
	a, _ := strconv.Atoi(m[4])
	b, _ := strconv.Atoi(m[6])
	if m[2] != "" {
		return Command{ParseOpCode(m[1]), ParseModifier(m[2]), ParseAddressingMode(m[3]), a, ParseAddressingMode(m[5]), b}, true
	}
	return Command{ParseOpCode(m[1]), opcodeStandard(m[1], m[3], m[5]), ParseAddressingMode(m[3]), a, ParseAddressingMode(m[5]), b}, true
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
