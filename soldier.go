package generals

import (
	"math"
	"fmt"
	"strings"
	"github.com/Arrow/display"
	"image/color"
)

const (
	vel             = 2.5
	spacing float64 = 10
	gain    float64 = 0.25
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

type Soldier struct {
	P          *display.Particle
	Pt         Point
	PastPt     Point
	Dir        float64
	Current    Order
	Past       Order
	ByLeft     bool
	FlankRatio float64 // Define it as a number from 0 to 1. 0 is right, 1 is left.
	OnRight    *Soldier
	InFront    *Soldier
	OnLeft     *Soldier
	Adj        []*Soldier
	Name       string
}

func NewSoldier(d *display.Display, name string, pt Point, dir float64) *Soldier {
	s := new(Soldier)
	s.Name = name
	s.P = d.NewParticle(pt.X, pt.Y, 2, RankColor)
	s.Pt = pt
	s.PastPt = pt
	s.Dir = dir
	s.Adj = make([]*Soldier, 4)
	return s
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
		str = append(str, fmt.Sprintf("R: %v; \n", s.Adj[1].GetName()))
	} else {
		str = append(str, "R:    nil; \n")
	}
	return strings.Join(str, "")
}

func (s *Soldier) Revert() {
	s.Current = s.Past
	s.Past = Halt
}

func (s *Soldier) Orders(o Order) {
	s.Past = s.Current
	s.Current = o
}

func (s *Soldier) Position() Point {
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
		return
	}
	if s.Adj[1] == nil {
		s.P.ChangeColor(RightColor)
		return
	}
	if s.Adj[3] == nil {
		s.P.ChangeColor(LeftColor)
		return
	}
	s.P.ChangeColor(color.Black)
}

func (s *Soldier) refPoint() Point {
	var pt, tmp Point
	var dir float64
	if s.Adj[0] != nil {
		dir = rotate(s.Dir, math.Pi)
		tmp = s.Adj[0].Position()
		pt.X += s.Pt.X-tmp.X-spacing*math.Cos(dir)
		pt.Y += s.Pt.Y-tmp.Y-spacing*math.Sin(dir)
	}
	if s.Adj[3] != nil && s.ByLeft {
		dir = rotate(s.Dir, -0.5*math.Pi)
		tmp = s.Adj[3].Position()
		pt.X += s.Pt.X-tmp.X-spacing*math.Cos(dir)
		pt.Y += s.Pt.Y-tmp.Y-spacing*math.Sin(dir)
	}
	if s.Adj[1] != nil && !s.ByLeft {
		dir = rotate(s.Dir, 0.5*math.Pi)
		tmp = s.Adj[1].Position()
		pt.X += s.Pt.X-tmp.X-spacing*math.Cos(dir)
		pt.Y += s.Pt.Y-tmp.Y-spacing*math.Sin(dir)
	}
	return pt
}

func (s *Soldier) Step() {
	ref := s.refPoint()
	switch s.Current {
	case Halt:
	case ForwardMarch:
		s.Pt.X += vel*math.Cos(s.Dir) - ref.X*gain
		s.Pt.Y += vel*math.Sin(s.Dir) - ref.Y*gain
		s.P.Move(s.Pt.X, s.Pt.Y)
	case LeftWheel:
	case RightWheel:
	case LeftTurn:
		s.Dir = leftTurn(s.Dir)
		s.Adj[0], s.Adj[1], s.Adj[2], s.Adj[3] = s.Adj[3], s.Adj[0], s.Adj[1], s.Adj[2]
		s.Color()
	case RightTurn:
		s.Dir = rightTurn(s.Dir)
		s.Adj[0], s.Adj[1], s.Adj[2], s.Adj[3] = s.Adj[1], s.Adj[2], s.Adj[3], s.Adj[0]
		s.Color()
	}
}
