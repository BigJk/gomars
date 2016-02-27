gomars
====

gomars is a CoreWar MARS written in go.

## example

```
m := gomars.CreateMars(8000, 8000, 80000, 200) // CoreSize, MaxProcess, Cycles, MaxLength

w1b, _ := ioutil.ReadFile("warrior1.red")
w2b, _ := ioutil.ReadFile("warrior2.red")

w1 := gomars.ParseWarrior(string(w1b))
w2 := gomars.ParseWarrior(string(w2b))

m.AddWarrior(w1)
m.AddWarrior(w2)

r := m.Run(1000)

fmt.Println("Wins #1", r[0])
fmt.Println("Wins #2", r[1])
fmt.Println("Ties", r[-1])
```

## todo

+ Adding org and end
+ Adding a parser (currently only "compiled" warriors will work)
+ Finishing MultiMARS for multi-core support
