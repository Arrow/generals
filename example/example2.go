package main

import (
	"fmt"
	"github.com/Arrow/display"
	"github.com/Arrow/generals"
	"log"
	"time"
)

const (
	width, height = 500, 500
	border        = 5
	heading       = 20
)

type Frame struct{}
type Timing struct{}

type OrderTiming struct {
	o generals.Order
	t time.Duration
}

func OrderKey(d *display.Display, ch chan generals.Order, o generals.Order, key string) {
	err := d.NewKeyBinding(func() { ch <- o }, key)
	if err != nil {
		log.Fatal(err)
	}
}

func KeyOrders(d *display.Display) chan generals.Order {
	ch := make(chan generals.Order)
	OrderKey(d, ch, generals.Halt, "h")
	OrderKey(d, ch, generals.ForwardMarch, "f")
	OrderKey(d, ch, generals.LeftTurn, "l")
	OrderKey(d, ch, generals.RightTurn, "r")
	return ch
}

func CallOrders(orders []OrderTiming) chan generals.Order {
	ch := make(chan generals.Order)

	go func() {
		for _, ord := range orders {
			time.Sleep(ord.t)
			ch <- ord.o
		}
	}()
	return ch
}

func main() {
	s1 := new(generals.Soldier)
	s1.Pt.X = 100
	s1.Pt.Y = 200

	s2 := new(generals.Soldier)
	s2.Pt.X = 95
	s2.Pt.Y = 200
	s2.InFront = s1

	d, err := display.NewDisplay(width, height, border, heading, "Soldier")
	if err != nil {
		log.Fatal(err)
	}

	s1.P = d.NewParticle(s1.Pt.X, s1.Pt.Y, 2, generals.FrontColor)
	s2.P = d.NewParticle(s2.Pt.X, s2.Pt.Y, 2, generals.RankColor)

	tick := time.Tick(500 * time.Millisecond)
	timer := time.Tick(time.Second)
	timerEnd := time.Tick(time.Minute)
	orders := KeyOrders(d)

	ctr := 0
	fps := 0
	for {
		<-tick
		s1.Step()
		s2.Step()
		s1.Update()
		s2.Update()
		d.Frame()
		ctr++
		select {
		case <-timer:
			fps = ctr
			d.SetHeadingText(fmt.Sprint("FPS: ", fps))
			ctr = 0
		case <-timerEnd:
			return
		case o := <-orders:
			s1.Orders(o)
			s2.Orders(o)
		default:
		}
	}
}
