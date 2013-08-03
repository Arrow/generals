package generals

import (
	"fmt"
	//"math"
	"github.com/Arrow/display"
	"github.com/golang/glog"
)

type Formation int

const (
	SingleCol Formation = iota
	TwoCol
	ThreeCol
	FourCol
	FiveCol
	SixCol
	SingleRow
	TwoRow
	ThreeRow
	FourRow
)

func (f Formation) RowCols(nSoldiers int) (rows, cols int) {
	switch f {
	case SingleRow, TwoRow, ThreeRow, FourRow:
		rows = int(f) - 5
		cols = nSoldiers / rows
		if nSoldiers%rows != 0 {
			cols += 1
		}
	default:
		cols = int(f) + 1
		rows = nSoldiers / cols
		if nSoldiers%cols != 0 {
			rows += 1
		}
	}
	return rows, cols
}

type Company struct {
	s       []*Soldier
	alignBy *Soldier
	f       Formation
}

func NewCompany(d *display.Display, pt complex128, dir complex128, nSoldiers int, f Formation) (c *Company) {
	c = new(Company)
	c.s = make([]*Soldier, nSoldiers)
	c.f = f
	_, cols := c.f.RowCols(len(c.s))
	for i, _ := range c.s {
		if i == 0 {
			c.s[i] = NewSoldier(d, fmt.Sprintf("Sol %v ", i), pt, dir)
			c.alignBy = c.s[0]
			glog.Info(c.s[i].GetName(), c.s[i].Pt)
		} else {
			c.s[i] = NewRandSoldier(d, fmt.Sprintf("Sol %v ", i))
			c.s[0].AddToForm(c.s[i], 1, cols)
			c.s[i].Pt = c.s[i].refPt()
			c.s[i].PastPt = c.s[i].Pt
			c.s[i].Dir = dir
			glog.Info(c.s[i].GetName(), c.s[i].Pt)
		}
	}
	for _, s := range c.s {
		s.Color()
		s.C = c
	}
	return c
}

func (c *Company) Step() {
	for _, s := range c.s {
		s.Step()
	}
}

func (c *Company) Update() {
	for _, s := range c.s {
		s.Update()
	}
}

func (c *Company) Orders(o Order) {
	glog.Info("Order", o)
	ord := o
	switch ord {
	case FormRow:
		c.f = SingleRow
		ord = Reform
	case FormTwoRow:
		c.f = TwoRow
		ord = Reform
	case FormThreeRow:
		c.f = ThreeRow
		ord = Reform
	case FormFourRow:
		c.f = FourRow
		ord = Reform
	case FormCol:
		c.f = SingleCol
		ord = Reform
	case FormTwoCol:
		c.f = TwoCol
		ord = Reform
	case FormThreeCol:
		c.f = ThreeCol
		ord = Reform
	case FormFourCol:
		c.f = FourCol
		ord = Reform
	case FormFiveCol:
		c.f = FiveCol
		ord = Reform
	}
	if ord == Reform {
		c.alignBy.nilAdj()
	}
	_, cols := c.f.RowCols(len(c.s))
	for _, s := range c.s {
		switch ord {
		case LeftTurn:
			s.lTurn()
		case RightTurn:
			s.rTurn()
		}
		if ord == Reform && s != c.alignBy {
			s.nilAdj()
			c.alignBy.AddToForm(s, 1, cols)
			s.Color()
		}
		s.Orders(ord)
	}
	if ord == LeftTurn || ord == RightTurn {
		c.alignBy = c.alignBy.findAlign()
	}
}

func (c *Company) PrintFormation() {
	fmt.Println("Formation Printed")
	for _, s := range c.s {
		glog.Info(s)
	}
}
