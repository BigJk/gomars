package gomars

func (c *Core) mov(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case a:
		c.Memory[bAddr].A = c.Memory[aAddr].A
	case b:
		c.Memory[bAddr].B = c.Memory[aAddr].B
	case ab:
		c.Memory[bAddr].B = c.Memory[aAddr].A
	case ba:
		c.Memory[bAddr].A = c.Memory[aAddr].B
	case x:
		c.Memory[bAddr].B = c.Memory[aAddr].A
		c.Memory[bAddr].A = c.Memory[aAddr].B
	case f:
		c.Memory[bAddr].B = c.Memory[aAddr].B
		c.Memory[bAddr].A = c.Memory[aAddr].A
	case i:
		c.Memory[bAddr] = c.Memory[aAddr].Clone()
	}
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
}

func (c *Core) add(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case a:
		c.Memory[bAddr].A = c.Memory[bAddr].A + c.Memory[aAddr].A
	case b:
		c.Memory[bAddr].B = c.Memory[bAddr].B + c.Memory[aAddr].B
	case ab:
		c.Memory[bAddr].B = c.Memory[bAddr].B + c.Memory[aAddr].A
	case ba:
		c.Memory[bAddr].A = c.Memory[bAddr].A + c.Memory[aAddr].B
	case x:
		c.Memory[bAddr].B = c.Memory[bAddr].B + c.Memory[aAddr].A
		c.Memory[bAddr].A = c.Memory[bAddr].A + c.Memory[aAddr].B
	case f, i:
		c.Memory[bAddr].B = c.Memory[bAddr].B + c.Memory[aAddr].B
		c.Memory[bAddr].A = c.Memory[bAddr].A + c.Memory[aAddr].A
	}
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
}

func (c *Core) sub(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case a:
		c.Memory[bAddr].A = c.Memory[bAddr].A - c.Memory[aAddr].A
	case b:
		c.Memory[bAddr].B = c.Memory[bAddr].B - c.Memory[aAddr].B
	case ab:
		c.Memory[bAddr].B = c.Memory[bAddr].B - c.Memory[aAddr].A
	case ba:
		c.Memory[bAddr].A = c.Memory[bAddr].A - c.Memory[aAddr].B
	case x:
		c.Memory[bAddr].B = c.Memory[bAddr].B - c.Memory[aAddr].A
		c.Memory[bAddr].A = c.Memory[bAddr].A - c.Memory[aAddr].B
	case f, i:
		c.Memory[bAddr].B = c.Memory[bAddr].B - c.Memory[aAddr].B
		c.Memory[bAddr].A = c.Memory[bAddr].A - c.Memory[aAddr].A
	}
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
}

func (c *Core) mul(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case a:
		c.Memory[bAddr].A = c.Memory[bAddr].A * c.Memory[aAddr].A
	case b:
		c.Memory[bAddr].B = c.Memory[bAddr].B * c.Memory[aAddr].B
	case ab:
		c.Memory[bAddr].B = c.Memory[bAddr].B * c.Memory[aAddr].A
	case ba:
		c.Memory[bAddr].A = c.Memory[bAddr].A * c.Memory[aAddr].B
	case x:
		c.Memory[bAddr].B = c.Memory[bAddr].B * c.Memory[aAddr].A
		c.Memory[bAddr].A = c.Memory[bAddr].A * c.Memory[aAddr].B
	case f, i:
		c.Memory[bAddr].B = c.Memory[bAddr].B * c.Memory[aAddr].B
		c.Memory[bAddr].A = c.Memory[bAddr].A * c.Memory[aAddr].A
	}
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
}

func (c *Core) div(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case a:
		if c.Memory[aAddr].A == 0 {
			return
		}
		c.Memory[bAddr].A = c.Memory[bAddr].A * c.Memory[aAddr].A
	case b:
		if c.Memory[aAddr].B == 0 {
			return
		}
		c.Memory[bAddr].B = c.Memory[bAddr].B * c.Memory[aAddr].B
	case ab:
		if c.Memory[aAddr].A == 0 {
			return
		}
		c.Memory[bAddr].B = c.Memory[bAddr].B * c.Memory[aAddr].A
	case ba:
		if c.Memory[aAddr].B == 0 {
			return
		}
		c.Memory[bAddr].A = c.Memory[bAddr].A * c.Memory[aAddr].B
	case x:
		if c.Memory[aAddr].A == 0 || c.Memory[aAddr].B == 0 {
			return
		}
		c.Memory[bAddr].B = c.Memory[bAddr].B * c.Memory[aAddr].A
		c.Memory[bAddr].A = c.Memory[bAddr].A * c.Memory[aAddr].B
	case f, i:
		if c.Memory[aAddr].A == 0 || c.Memory[aAddr].B == 0 {
			return
		}
		c.Memory[bAddr].B = c.Memory[bAddr].B * c.Memory[aAddr].B
		c.Memory[bAddr].A = c.Memory[bAddr].A * c.Memory[aAddr].A
	}
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
}

func (c *Core) mod(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case a:
		if c.Memory[aAddr].A == 0 {
			return
		}
		c.Memory[bAddr].A = c.Memory[bAddr].A % c.Memory[aAddr].A
	case b:
		if c.Memory[aAddr].B == 0 {
			return
		}
		c.Memory[bAddr].B = c.Memory[bAddr].B % c.Memory[aAddr].B
	case ab:
		if c.Memory[aAddr].A == 0 {
			return
		}
		c.Memory[bAddr].B = c.Memory[bAddr].B % c.Memory[aAddr].A
	case ba:
		if c.Memory[aAddr].B == 0 {
			return
		}
		c.Memory[bAddr].A = c.Memory[bAddr].A % c.Memory[aAddr].B
	case x:
		if c.Memory[aAddr].A == 0 || c.Memory[aAddr].B == 0 {
			return
		}
		c.Memory[bAddr].B = c.Memory[bAddr].B % c.Memory[aAddr].A
		c.Memory[bAddr].A = c.Memory[bAddr].A % c.Memory[aAddr].B
	case f, i:
		if c.Memory[aAddr].A == 0 || c.Memory[aAddr].B == 0 {
			return
		}
		c.Memory[bAddr].B = c.Memory[bAddr].B % c.Memory[aAddr].B
		c.Memory[bAddr].A = c.Memory[bAddr].A % c.Memory[aAddr].A
	}
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
}

func (c *Core) jmp(aAddr int, baseAddress int, w *CoreWarrior) {
	w.QueueTask(c.NormalizeAddress(aAddr))
}

func (c *Core) jmz(aAddr, bAddr int, baseAddress int, w *CoreWarrior) {
	if c.Memory[bAddr].B == 0 {
		w.QueueTask(c.NormalizeAddress(aAddr))
	} else {
		w.QueueTask(c.NormalizeAddress(baseAddress + 1))
	}
}

func (c *Core) jmn(aAddr, bAddr int, baseAddress int, w *CoreWarrior) {
	if c.Memory[bAddr].B != 0 {
		w.QueueTask(c.NormalizeAddress(aAddr))
	} else {
		w.QueueTask(c.NormalizeAddress(baseAddress + 1))
	}
}

func (c *Core) djn(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case a, ba:
		c.Memory[bAddr].A--
		if c.Memory[bAddr].A != 0 {
			w.QueueTask(c.NormalizeAddress(aAddr))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case b, ab:
		c.Memory[bAddr].B--
		if c.Memory[bAddr].B != 0 {
			w.QueueTask(c.NormalizeAddress(aAddr))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case f, x, i:
		c.Memory[bAddr].B--
		c.Memory[bAddr].A--
		if c.Memory[bAddr].B != 0 {
			w.QueueTask(c.NormalizeAddress(aAddr))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	}
}

func (c *Core) spl(aAddr int, baseAddress int, w *CoreWarrior) {
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
	w.QueueTask(c.NormalizeAddress(aAddr))
}

func (c *Core) seq(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case i:
		if c.Memory[aAddr].Equal(&c.Memory[bAddr]) {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case a:
		if c.Memory[aAddr].A == c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case b:
		if c.Memory[aAddr].B == c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case ab:
		if c.Memory[aAddr].A == c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case ba:
		if c.Memory[aAddr].B == c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case f:
		if c.Memory[aAddr].A == c.Memory[bAddr].A && c.Memory[aAddr].B == c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case x:
		if c.Memory[aAddr].A == c.Memory[bAddr].B && c.Memory[aAddr].B == c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	}
}

func (c *Core) sne(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case i:
		if !c.Memory[aAddr].Equal(&c.Memory[bAddr]) {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case a:
		if c.Memory[aAddr].A != c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case b:
		if c.Memory[aAddr].B != c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case ab:
		if c.Memory[aAddr].A != c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case ba:
		if c.Memory[aAddr].B != c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case f:
		if c.Memory[aAddr].A != c.Memory[bAddr].A && c.Memory[aAddr].B != c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case x:
		if c.Memory[aAddr].A != c.Memory[bAddr].B && c.Memory[aAddr].B != c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	}
}

func (c *Core) slt(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	switch m {
	case a:
		if c.Memory[aAddr].A < c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case b:
		if c.Memory[aAddr].B < c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case ab:
		if c.Memory[aAddr].A < c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case ba:
		if c.Memory[aAddr].B < c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case i, f:
		if c.Memory[aAddr].A < c.Memory[bAddr].A && c.Memory[aAddr].B < c.Memory[bAddr].B {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	case x:
		if c.Memory[aAddr].A < c.Memory[bAddr].B && c.Memory[aAddr].B < c.Memory[bAddr].A {
			w.QueueTask(c.NormalizeAddress(baseAddress + 2))
		} else {
			w.QueueTask(c.NormalizeAddress(baseAddress + 1))
		}
	}
}

func (c *Core) stp(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
	switch m {
	case a:
		p := c.NormalizePSpaceAddress(c.Memory[bAddr].A)
		if p == 0 {
			return
		}
		c.PSpace[w.ID][p] = c.Memory[aAddr].A
	case f, x, i, b:
		p := c.NormalizePSpaceAddress(c.Memory[bAddr].B)
		if p == 0 {
			return
		}
		c.PSpace[w.ID][p] = c.Memory[aAddr].B
	case ab:
		p := c.NormalizePSpaceAddress(c.Memory[bAddr].B)
		if p == 0 {
			return
		}
		c.PSpace[w.ID][p] = c.Memory[aAddr].A
	case ba:
		p := c.NormalizePSpaceAddress(c.Memory[bAddr].A)
		if p == 0 {
			return
		}
		c.PSpace[w.ID][p] = c.Memory[aAddr].B
	}
}

func (c *Core) ldp(aAddr, bAddr int, m Modifier, baseAddress int, w *CoreWarrior) {
	w.QueueTask(c.NormalizeAddress(baseAddress + 1))
	switch m {
	case a:
		c.Memory[bAddr].A = c.PSpace[w.ID][c.NormalizePSpaceAddress(c.Memory[aAddr].A)]
	case f, x, i, b:
		c.Memory[bAddr].B = c.PSpace[w.ID][c.NormalizePSpaceAddress(c.Memory[aAddr].B)]
	case ab:
		c.Memory[bAddr].B = c.PSpace[w.ID][c.NormalizePSpaceAddress(c.Memory[aAddr].A)]
	case ba:
		c.Memory[bAddr].A = c.PSpace[w.ID][c.NormalizePSpaceAddress(c.Memory[aAddr].B)]
	}
}
