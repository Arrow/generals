package generals

import (
	//"fmt"
	"math"
)

// RegPlacement defines the placement of a regiment.
type RegPlacement struct {
	Dir float64 // in Radians from 0 to 2 Pi
	//Rows int
	RowPath []Point
	//Columns int
	ColPath []Point
}

func ColumnRP(x, y, dir float64, cols, soldiers int) *RegPlacement {
	rp := new(RegPlacement)
	rp.Dir = dir
	rows := soldiers/cols + 1
	rp.RowPath = make([]Point, rows)
	rp.ColPath = make([]Point, cols)
	dr := rotate(dir, math.Pi)
	for i := 0; i < rows; i++ {
		rp.RowPath[i].X = x + float64(i)*spacing*math.Cos(dr)
		rp.RowPath[i].Y = y - float64(i)*spacing*math.Sin(dr)
	}
	dc := rotate(dir, 0.5*math.Pi)
	for i := 0; i < cols; i++ {
		rp.ColPath[i].X = x + float64(i)*spacing*math.Cos(dc)
		rp.ColPath[i].Y = y - float64(i)*spacing*math.Sin(dc)
	}
	return rp
}

func RowRP(x, y, dir float64, rows, soldiers int) *RegPlacement {
	rp := new(RegPlacement)
	rp.Dir = dir
	cols := soldiers/rows + 1
	rp.RowPath = make([]Point, rows)
	rp.ColPath = make([]Point, cols)
	dr := rotate(dir, math.Pi)
	for i := 0; i < rows; i++ {
		rp.RowPath[i].X = x + float64(i)*spacing*math.Cos(dr)
		rp.RowPath[i].Y = y - float64(i)*spacing*math.Sin(dr)
	}
	dc := rotate(dir, 0.5*math.Pi)
	for i := 0; i < cols; i++ {
		rp.ColPath[i].X = x + float64(i)*spacing*math.Cos(dc)
		rp.ColPath[i].Y = y - float64(i)*spacing*math.Sin(dc)
	}
	return rp
}

// Regiment defines the behaviour
type Regiment struct {
	name           string
	rp             *RegPlacement
	numSoldiers    int
	musterStrength int
	current        Order
	past           Order
}

func NewRegiment(name string, rp *RegPlacement, muster int) *Regiment {
	reg := new(Regiment)
	reg.name = name
	reg.rp = rp
	reg.numSoldiers = muster
	reg.musterStrength = muster
	reg.current = Halt
	reg.past = Halt
	return reg
}

func (r *Regiment) Revert() {
	r.current = r.past
	r.past = Halt
}

func (r *Regiment) Orders(o Order) {
	r.past = r.current
	r.current = o
}

func (r *Regiment) Position() *RegPlacement {
	return r.rp
}

func (r *Regiment) Step() {
	switch r.current {
	case Halt:
	case ForwardMarch:
		for i := 0; i < len(r.rp.ColPath); i++ {
			r.rp.ColPath[i].X += vel * math.Cos(r.rp.Dir)
			r.rp.ColPath[i].Y -= vel * math.Sin(r.rp.Dir)
		}
		x, y := r.rp.ColPath[0].X, r.rp.ColPath[0].Y
		for i := 0; i < len(r.rp.RowPath); i++ {
			facX, facY := norm(x-r.rp.RowPath[i].X, y-r.rp.RowPath[i].Y)
			x, y = r.rp.RowPath[i].X, r.rp.RowPath[i].Y
			r.rp.RowPath[i].X += vel * facX
			r.rp.RowPath[i].Y += vel * facY
		}
	case LeftWheel:
	case RightWheel:
	case LeftTurn:
		r.rp.Dir = leftTurn(r.rp.Dir)
		lr, lc := len(r.rp.ColPath), len(r.rp.RowPath)
		rows := make([]Point, lr)
		cols := make([]Point, lc)
		colOffset := float64((lr - 1)) * spacing
		for i := 0; i < lc; i++ {
			cols[i] = Point{
				r.rp.RowPath[i].X + colOffset*math.Cos(r.rp.Dir),
				r.rp.RowPath[i].Y - colOffset*math.Sin(r.rp.Dir),
			}
		}
		for i := 0; i < lr; i++ {
			rows[i] = r.rp.ColPath[lr-1-i]
		}
		r.rp.RowPath = rows
		r.rp.ColPath = cols
		r.Revert()
	case RightTurn:
		r.rp.Dir = rightTurn(r.rp.Dir)
		lr, lc := len(r.rp.ColPath), len(r.rp.RowPath)
		rows := make([]Point, lr)
		cols := make([]Point, lc)

		rowOffset := float64((lc - 1)) * spacing
		for i := 0; i < lr; i++ {
			rows[i] = Point{
				r.rp.ColPath[i].X + rowOffset*math.Cos(rightTurn(r.rp.Dir)),
				r.rp.ColPath[i].Y - rowOffset*math.Sin(rightTurn(r.rp.Dir)),
			}
		}
		for i := 0; i < lc; i++ {
			cols[i] = r.rp.RowPath[lc-1-i]
		}

		r.rp.RowPath = rows
		r.rp.ColPath = cols
		r.Revert()
	}
}
