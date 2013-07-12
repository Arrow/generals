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
	OrderKey(d, ch, generals.RightWheel, "w")
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
	d, err := display.NewDisplay(width, height, border, heading, "Soldier")
	if err != nil {
		log.Fatal(err)
	}
	
	c := generals.NewCompany(d, generals.Point{100, 200}, 0, 25, generals.TwoRow)
	
	tick := time.Tick(500 * time.Millisecond)
	timer := time.Tick(time.Second)
	timerEnd := time.Tick(60 * time.Second)
	orders := KeyOrders(d)
	
	ctr := 0
	fps := 0
	for {
		<-tick
		c.Step()
		c.Update()
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
			if o == generals.RightWheel {
				return
			}
			c.Orders(o)
		default:
		}
	}
}
