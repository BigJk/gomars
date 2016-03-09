package gomars

import (
	"fmt"
	"testing"
)

func BenchmarkMars(b *testing.B) {
	m := CreateMars(800, 800, 8000, 20, 20)
	m.AddWarrior(m.ParseWarriorFromFile("test/1.red"))
	m.AddWarrior(m.ParseWarriorFromFile("test/2.red"))
	r := m.Run(2000)
	fmt.Println(r)
}
