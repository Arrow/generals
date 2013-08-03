package generals

import (
	"fmt"
	"github.com/Arrow/display"
	"github.com/golang/glog"
	"image/color"
	"math"
	"math/cmplx"
	"math/rand"
	"strings"
)

const (
	vel     float64 = 2.5
	velMax  float64 = 4.0
	spacing float64 = 7.5
	gain    float64 = 0.35
	gainRef float64 = 0.50
	velRand float64 = 0.20
)

var (
	FrontColor = color.RGBA{255, 0, 0, 255}
	LeftColor  = color.RGBA{0, 255, 0, 255}
	RightColor = color.RGBA{0, 0, 255, 255}
	RankColor  = color.Black
)

type Point struct {
	X float64
	Y float64
}

func norm(x float64, y float64) (float64, float64) {
	fac := 1 / math.Sqrt(x*x+y*y)
	return x * fac, y * fac
}

func leftTurn(dir float64) float64 {
	tmp := dir + math.Pi/2
	if tmp > 2*math.Pi {
		tmp -= 2 * math.Pi
	}
	return tmp
}

func rightTurn(dir float64) float64 {
	tmp := dir - math.Pi/2
	if tmp < 0 {
		tmp += 2 * math.Pi
	}
	return tmp
}

func rotate(dir, angle float64) float64 {
	tmp := dir + angle
	if tmp > 2*math.Pi {
		tmp -= 2 * math.Pi
	} else if tmp < 0 {
		tmp += 2 * math.Pi
	}
	return tmp
}

type FormPos struct {
	Row int
	Col int
}

type Soldier struct {
	P          *display.Particle
	Pt         complex128
	PastPt     complex128
	Dir        complex128
	Ps         FormPos
	Current    Order
	Past       Order
	ByLeft     bool
	FlankRatio float64 // Define it as a number from 0 to 1. 0 is right, 1 is left.
	Adj        []*Soldier
	Name       string
	C          *Company
}

func NewSoldier(d *display.Display, name string, pt complex128, dir complex128) *Soldier {
	s := new(Soldier)
	s.Name = name
	s.P = d.NewParticle(real(pt), imag(pt), 1, RankColor)
	s.Pt = pt
	s.PastPt = pt
	s.Dir = dir
	s.Adj = make([]*Soldier, 4)
	return s
}

func NewRandSoldier(d *display.Display, name string) *Soldier {
	s := new(Soldier)
	s.Name = name
	s.Pt = complex(rand.Float64(), rand.Float64())
	s.P = d.NewParticle(real(s.Pt), imag(s.Pt), 1, RankColor)
	s.PastPt = s.Pt
	s.Adj = make([]*Soldier, 4)
	return s
}

func (s *Soldier) nilAdj() {
	s.Adj[0], s.Adj[1], s.Adj[2], s.Adj[3] =
		nil, nil, nil, nil
}

func (s *Soldier) findAlign() *Soldier {
	if s.Adj[0] == nil {
		if s.Adj[1] == nil {
			return s
		}
		return s.Adj[1].findAlign()
	}
	return s.Adj[0].findAlign()
}

func (s *Soldier) AddToForm(sc *Soldier, ctr, cols int) (added bool) {
	glog.V(2).Info(s, sc, ctr, cols)
	if ctr == cols {
		return false
	}
	if s.Adj[3] == nil {
		s.Adj[3] = sc
		sc.Adj[1] = s
		if s.Adj[0] != nil {
			if s.Adj[0].Adj[3] != nil {
				s.Adj[0].Adj[3].Adj[2] = sc
				sc.Adj[0] = s.Adj[0].Adj[3]
			}
			return true
		}
		return true
	}
	if s.Adj[3].AddToForm(sc, ctr+1, cols) {
		return true
	}
	if s.Adj[1] == nil {
		if s.Adj[2] == nil {
			s.Adj[2] = sc
			sc.Adj[0] = s
			return true
		}
		s.Adj[2].AddToForm(sc, ctr, cols)
		return true
	}
	return false
}

func (s *Soldier) GetName() string {
	return s.Name
}

func (s *Soldier) String() string {
	str := make([]string, 0)
	str = append(str, fmt.Sprintf("Name: %v; ", s.Name))
	if s.Adj[3] != nil {
		str = append(str, fmt.Sprintf("L: %v; ", s.Adj[3].GetName()))
	} else {
		str = append(str, "L:    nil; ")
	}
	if s.Adj[0] != nil {
		str = append(str, fmt.Sprintf("F: %v; ", s.Adj[0].GetName()))
	} else {
		str = append(str, "F:    nil; ")
	}
	if s.Adj[1] != nil {
		str = append(str, fmt.Sprintf("R: %v; ", s.Adj[1].GetName()))
	} else {
		str = append(str, "R:    nil; ")
	}
	return strings.Join(str, "")
}

func (s *Soldier) Revert() {
	s.Current = s.Past
	s.Past = Halt
	s.Color()
}

func (s *Soldier) Orders(o Order) {
	s.Past = s.Current
	s.Current = o
}

func (s *Soldier) Position() complex128 {
	return s.PastPt
}

func (s *Soldier) Update() {
	s.PastPt = s.Pt
	switch s.Current {
	case LeftTurn, RightTurn:
		s.Revert()
	}
}

func (s *Soldier) Color() {
	if s.Adj[0] == nil {
		s.P.ChangeColor(FrontColor)
		s.P.Move(real(s.Pt), imag(s.Pt))
		return
	}
	if s.Adj[1] == nil {
		s.P.ChangeColor(RightColor)
		s.P.Move(real(s.Pt), imag(s.Pt))
		return
	}
	if s.Adj[3] == nil {
		s.P.ChangeColor(LeftColor)
		s.P.Move(real(s.Pt), imag(s.Pt))
		return
	}
	s.P.ChangeColor(color.Black)
	s.P.Move(real(s.Pt), imag(s.Pt))
}

/*
func (s *Soldier) Row() int {
	if s.Adj[0] == nil {
		return 0
	}
	return s.Adj[0].Row() + 1
}
*/
func (s *Soldier) Col() int {
	if s.Adj[1] == nil {
		return 1
	}
	return s.Adj[1].Col() + 1
}

var (
	left  complex128 = cmplx.Rect(1, 0.5*math.Pi)
	right complex128 = cmplx.Rect(1, -0.5*math.Pi)
)

func (s *Soldier) refPoint() complex128 {
	var pt, tmp, dir complex128
	if s.Adj[0] != nil {
		dir = s.Dir * left * left
		tmp = s.Adj[0].Position()
		pt += s.Pt - tmp - complex(spacing, 0)*dir
	}
	if s.Adj[3] != nil && s.ByLeft {
		dir = s.Dir * right
		tmp = s.Adj[3].Position()
		pt += s.Pt - tmp - complex(spacing, 0)*dir
	}
	if s.Adj[1] != nil && !s.ByLeft {
		dir = s.Dir * left
		tmp = s.Adj[1].Position()
		pt += s.Pt - tmp - complex(spacing, 0)*dir
	}
	return pt
}

func (s *Soldier) refPt() complex128 {
	var pt, tmp, dir complex128
	var ctr int
	if s.Adj[0] != nil {
		ctr++
		dir = s.Adj[0].Dir
		tmp = s.Adj[0].Position()
		pt += tmp - complex(spacing, 0)*dir
	}
	if s.Adj[3] != nil && s.ByLeft {
		ctr++
		dir = s.Adj[3].Dir * left
		tmp = s.Adj[3].Position()
		pt += tmp - complex(spacing, 0)*dir
	}
	if s.Adj[1] != nil && !s.ByLeft {
		ctr++
		dir = s.Adj[1].Dir * right
		tmp = s.Adj[1].Position()
		pt += tmp - complex(spacing, 0)*dir
	}
	glog.Info(pt, ctr)
	return pt / complex(float64(ctr), 0)
}

func (s *Soldier) Step() {
	ref := s.refPoint()
	switch s.Current {
	case Halt:
	case ForwardMarch:
		s.Pt += complex(vel, 0)*s.Dir - ref*complex(gain, 0) +
			complex(velRand, 0)*complex(rand.Float64(), rand.Float64())
		s.P.Move(real(s.Pt), imag(s.Pt))
	case LeftWheel:
		v := velMax
		if s.Adj[0] == nil {
			_, cols := s.C.f.RowCols(len(s.C.s))
			v *= float64(cols-s.Col()+1) / float64(cols+1)
			th := velMax / (2 * math.Pi * float64(cols+1))
			s.Dir *= cmplx.Rect(1, th)
		} else {
			v = vel
		}
		s.Pt += complex(v, 0)*s.Dir - ref*complex(gain, 0) +
			complex(velRand, 0)*complex(rand.Float64(), rand.Float64())
		s.P.Move(real(s.Pt), imag(s.Pt))
	case RightWheel:
		v := velMax
		if s.Adj[0] == nil {
			_, cols := s.C.f.RowCols(len(s.C.s))
			v *= float64(s.Col()+1) / float64(cols+1)
			th := velMax / (2 * math.Pi * float64(cols+1))
			s.Dir *= cmplx.Rect(1, -th)
		} else {
			v = vel
		}
		s.Pt += complex(v, 0)*s.Dir - ref*complex(gain, 0) +
			complex(velRand, 0)*complex(rand.Float64(), rand.Float64())
		s.P.Move(real(s.Pt), imag(s.Pt))
		//	case LeftTurn, RightTurn:
		//		s.Color()
	case Reform:
		s.Pt -= ref*complex(gainRef, 0) -
			complex(velRand, 0)*complex(rand.Float64(), rand.Float64())
		s.P.Move(real(s.Pt), imag(s.Pt))
		if cmplx.Abs(s.Pt-ref) < 0.1 {
			s.Orders(Halt)
		}
	}
	if s.Adj[0] != nil {
		s.Dir = cmplx.Rect(1, cmplx.Phase(s.Adj[0].Position()-s.Position()))
	}
}

func (s *Soldier) lTurn() {
	s.Dir *= left
	s.Adj[0], s.Adj[1], s.Adj[2], s.Adj[3] =
		s.Adj[3], s.Adj[0], s.Adj[1], s.Adj[2]
}

func (s *Soldier) rTurn() {
	s.Dir *= right
	s.Adj[0], s.Adj[1], s.Adj[2], s.Adj[3] =
		s.Adj[1], s.Adj[2], s.Adj[3], s.Adj[0]
}
