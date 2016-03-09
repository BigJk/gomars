package gomars

import (
	"io/ioutil"
	"strings"
	"testing"
)

var testWarrior = [][]string{
	[]string{"1.red", "1.rc"},
	[]string{"2.red", "2.rc"},
	[]string{"3.red", "3.rc"},
	[]string{"4.red", "4.rc"},
	[]string{"5.red", "5.rc"},
	[]string{"6.red", "6.rc"},
	[]string{"7.red", "7.rc"}}

func TestParser(t *testing.T) {
	m := CreateMars(800, 800, 8000, 20, 20)

	for i := 0; i < len(testWarrior); i++ {
		w, _ := ioutil.ReadFile("test/" + testWarrior[i][0])
		wrc, _ := ioutil.ReadFile("test/" + testWarrior[i][1])
		wout := m.ParseWarrior(string(w))

		wrcLines := singleLine.FindAllStringSubmatch(strings.ToLower(string(wrc)), -1)
		woutLines := singleLine.FindAllStringSubmatch(wout.LoadFile(), -1)

		if len(wrcLines) != len(woutLines) {
			t.Error("Parsing test failed with", testWarrior[i][0], "/", testWarrior[i][1], "Line count is different.")
		}

		for j := 0; j < len(wrcLines); j++ {
			if wrcLines[j][1] != woutLines[j][1] || wrcLines[j][2] != woutLines[j][2] || wrcLines[j][3] != woutLines[j][3] || wrcLines[j][4] != woutLines[j][4] || wrcLines[j][5] != woutLines[j][5] || wrcLines[j][6] != woutLines[j][6] {
				t.Error("Parsing test failed with", testWarrior[i][0], "/", testWarrior[i][1], "Line is different. Expected:", wrcLines[j][0], "got", woutLines[j][0])
			}
		}
	}

}
