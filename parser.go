package gomars

import (
	"fmt"
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
var isFormular = regexp.MustCompile("[0-9][+\\-*\\/%][0-9]")
var field = regexp.MustCompile("^([*#{}$<>@])(.*)")
var addressingModes = regexp.MustCompile("([a-z]{3}) ([*#{}$<>@]).*, ([*#{}$<>@])")
var isNumber = regexp.MustCompile("[0-9]+")

var opcodes = []string{"dat", "mov", "add", "sub", "mul", "div", "mod", "jmp", "jmz", "jmn", "djn", "spl", "seq", "sne", "cmp", "slt", "nop", "ldp", "stp"}
var modifier = []string{"*", "#", "{", "}", "$", ">", "<", "@"}
var env = eval.NewEnvironment()

// ParseToLoadFile ...
func (m *MARS) ParseToLoadFile(warrior string) string {
	predefined := make(map[string]int)
	predefined["CORESIZE"] = m.Coresize
	predefined["PSPACESIZE"] = len(m.Core.PSpace)
	predefined["MAXCYCLES"] = m.MaxCycles
	predefined["MAXPROCESSES"] = m.MaxProcess
	predefined["WARRIORS"] = len(m.warrior) // TODO: FIX!
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

	compiledWarrior := ""

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

		newLine := ""

		if len(split) >= 3 {
			newLine += split[0] + " " + padRight(compileField(split[1])) + ", " + compileField(split[2]) + "\n"
		} else {
			newLine += split[0] + " " + padRight(compileField(split[1])) + ", $0 \n"
		}

		if !strings.Contains(split[0], ".") {
			f := addressingModes.FindStringSubmatch(newLine)

			newLine = strings.Replace(newLine, f[1], f[1]+"."+opcodeStandard(f[1], f[2], f[3]), -1)
		}

		compiledWarrior += newLine

	}

	if org != "" {
		if !isNumber.MatchString(org) {
			org = fmt.Sprint(labels[org])
		}
	} else if end != "" {
		if !isNumber.MatchString(org) {
			org = fmt.Sprint(labels[end])
		}
	} else {
		org = "0"
	}

	return "ORG " + org + "\nPIN 0\n" + strings.ToUpper(compiledWarrior)
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
	return strings.Replace(s, "(-", "(0-", -1)
}

func replaceCurline(s string, i int) string {
	return strings.Replace(s, "CURLINE", strconv.Itoa(i), -1)
}

func compileField(s string) string {
	f := field.FindStringSubmatch(s)
	if len(f) == 0 {
		if len(isFormular.FindAllString(s, -1)) == 0 {
			return "$" + s
		}
		v, _ := env.Eval(fixForumlar(s))
		return "$" + fmt.Sprint(v)
	}
	if len(isFormular.FindAllString(s, -1)) == 0 {
		return s
	}
	v, _ := env.Eval(fixForumlar(f[2]))
	return f[1] + fmt.Sprint(v)
}

func isFor(s string) bool {
	return strings.ToLower(s) == "for"
}

func opcodeStandard(s string, a string, b string) string {
	o := strings.ToLower(s)

	if o == "dat" {
		return "f"
	}

	if (o == "mov" || o == "cmp") && a == "#" || (o == "add" || o == "sub" || o == "mul" || o == "div" || o == "mod" || o == "slt") && a == "#" {
		return "ab"
	} else if (o == "mov" || o == "cmp") && b == "#" || (o == "add" || o == "sub" || o == "mul" || o == "div" || o == "mod") && a == "#" {
		return "b"
	} else if o == "mov" || o == "cmp" {
		return "i"
	} else if o == "add" || o == "sub" || o == "mul" || o == "div" || o == "mod" {
		return "f"
	}

	return "b"
}

func padRight(s string) string {
	out := s
	for i := 0; i < 10-len(s); i++ {
		out += " "
	}
	return out
}
