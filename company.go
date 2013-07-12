package generals

import (
	//"fmt"
	"math"
	"github.com/Arrow/display"
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

type Company struct {
	s []*Soldier
	rows int
	cols int
}

func NewCompany(d *display.Display, pt Point, dir float64, nSoldiers int, f Formation) (c *Company) {
	c = new(Company)
	c.s = make([]*Soldier, nSoldiers)
	switch f {
	case SingleRow, TwoRow, ThreeRow, FourRow:
		c.rows = int(f) - 5
		c.cols = nSoldiers / c.rows
		if nSoldiers % c.rows != 0 {
			c.cols += 1
		}
	default:
		c.cols = int(f) + 1
		c.rows = nSoldiers / c.cols
		if nSoldiers % c.cols != 0 {
			c.rows += 1
		}
	}
	row := 0
	col := 0
	pti := pt
	for i, _ := range c.s {
		left := leftTurn(dir)
		back := leftTurn(left)
		pti.X = pt.X + 
			float64(col) * spacing * math.Cos(left) + 
			float64(row) * spacing * math.Cos(back)
		pti.Y = pt.Y + 
			float64(col) * spacing * math.Sin(left) + 
			float64(row) * spacing * math.Sin(back)
		c.s[i] = NewSoldier(d, pti, dir)
		col++
		if col == c.cols {
			col = 0
			row++
		}
	}
	//fmt.Println(c.cols, c.rows)
	for i, s := range c.s {
		if i >= c.cols {
			s.SetInFront(c.s[i-c.cols])
		}
		if i % c.cols != 0 {
			s.SetOnRight(c.s[i-1])
		}
		if i % c.cols != c.cols - 1 && i != len(c.s) - 1 {
			s.SetOnLeft(c.s[i+1])
		}
		col++
		if col == c.cols {
			col = 0
		}
	}
	//fmt.Println("Here")
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
	for _, s := range c.s {
		s.Orders(o)
		//fmt.Println(s)
	}
}
