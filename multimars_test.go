package gomars

import (
	"fmt"
	"testing"
)

func BenchmarkMultiMars(b *testing.B) {
	m := CreateMultiMars(800, 800, 8000, 20, 20, 6)
	m.AddWarrior(m.ParseWarriorFromFile("test/1.red"))
	m.AddWarrior(m.ParseWarriorFromFile("test/2.red"))
	r := m.Run(2000)
	fmt.Println(r)
}
