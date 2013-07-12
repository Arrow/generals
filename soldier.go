package generals

import (
	"math"
	//"fmt"
	//"image"
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
}

func NewSoldier(d *display.Display, pt Point, dir float64) *Soldier {
	s := new(Soldier)
	s.P = d.NewParticle(pt.X, pt.Y, 2, RankColor)
	s.Pt = pt
	s.PastPt = pt
	s.Dir = dir
	return s
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
}

func (s *Soldier) Color() color.Color {
	if s.InFront == nil {
		return FrontColor
	}
	if s.OnLeft == nil {
		return LeftColor
	}
	if s.OnRight == nil {
		return RightColor
	}
	return color.Black
}

func (s *Soldier) SetInFront(sRef *Soldier) {
	s.InFront = sRef
}

func (s *Soldier) SetOnRight(sRef *Soldier) {
	s.OnRight = sRef
}

func (s *Soldier) SetOnLeft(sRef *Soldier) {
	s.OnLeft = sRef
}

func (s *Soldier) refPoint() Point {
	var pt, tmp Point
	var dir float64
	if s.InFront != nil {
		dir = rotate(s.Dir, math.Pi)
		tmp = s.InFront.Position()
		pt.X += s.Pt.X-tmp.X-spacing*math.Cos(dir)
		pt.Y += s.Pt.Y-tmp.Y-spacing*math.Sin(dir)
	}
	if s.OnLeft != nil && s.ByLeft {
		dir = rotate(s.Dir, -0.5*math.Pi)
		tmp = s.OnLeft.Position()
		pt.X += s.Pt.X-tmp.X-spacing*math.Cos(dir)
		pt.Y += s.Pt.Y-tmp.Y-spacing*math.Sin(dir)
	}
	if s.OnRight != nil && !s.ByLeft {
		dir = rotate(s.Dir, 0.5*math.Pi)
		tmp = s.OnRight.Position()
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
		s.Revert()
	case RightTurn:
		s.Dir = rightTurn(s.Dir)
		s.Revert()
	}
}
