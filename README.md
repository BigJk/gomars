gomars
====

gomars is a CoreWar MARS written in go.

## example

```
m := gomars.CreateMars(8000, 8000, 80000, 200, 200) // CoreSize, MaxProcess, MaxCycles, MaxLength, MinDIstance

w1b, _ := ioutil.ReadFile("warrior1.red")
w2b, _ := ioutil.ReadFile("warrior2.red")

m.AddWarriorString(string(w1b))
m.AddWarriorString(string(w2b))

start := time.Now()
r := m.Run(100)
elapsed := time.Since(start)

fmt.Println("Wins #1", r[0])
fmt.Println("Wins #2", r[1])
fmt.Println("Ties", r[-1])
```

## todo

+ Fixing PSpace (implemented but not fully working yet)
+ Finishing MultiMARS for multi-core support
