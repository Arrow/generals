package main

import (
	"fmt"
	"github.com/Arrow/display"
	"github.com/Arrow/generals"
	//"log"
	"time"
	"github.com/golang/glog"
)

const (
	width, height = 500, 500
	border        = 5
	heading       = 20
)

type OrderTiming struct {
	o generals.Order
	t time.Duration
}

func OrderKey(d *display.Display, ch chan generals.Order, o generals.Order, key string) {
	err := d.NewKeyBinding(func() { ch <- o }, key)
	if err != nil {
		glog.Fatal(err)
	}
}

func KeyOrders(d *display.Display) chan generals.Order {
	ch := make(chan generals.Order)
	OrderKey(d, ch, generals.Halt, "s")
	OrderKey(d, ch, generals.ForwardMarch, "w")
	OrderKey(d, ch, generals.LeftTurn, "a")
	OrderKey(d, ch, generals.RightTurn, "d")
	OrderKey(d, ch, generals.LeftWheel, "q")
	OrderKey(d, ch, generals.RightWheel, "e")
	OrderKey(d, ch, generals.Quit, "o")
	OrderKey(d, ch, generals.PrintForm, "p")
	OrderKey(d, ch, generals.FormRow, "1")
	OrderKey(d, ch, generals.FormTwoRow, "2")
	OrderKey(d, ch, generals.FormThreeRow, "3")
	OrderKey(d, ch, generals.FormFourRow, "4")
	OrderKey(d, ch, generals.FormCol, "g")
	OrderKey(d, ch, generals.FormTwoCol, "h")
	OrderKey(d, ch, generals.FormThreeCol, "j")
	OrderKey(d, ch, generals.FormFourCol, "k")
	OrderKey(d, ch, generals.FormFiveCol, "l")
	glog.Infoln("Keys Logged")
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
	defer glog.Flush()
	
	d, err := display.NewDisplay(width, height, border, heading, "Soldier")
	if err != nil {
		glog.Fatal(err)
	}
	
	c := generals.NewCompany(d, complex(float64(100), float64(200)), 
		complex(float64(1),float64(0)), 24, generals.FourCol)
	
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
			//return
		case o := <-orders:
			if o == generals.Quit {
				return
			}
			if o == generals.PrintForm {
				c.PrintFormation()
				break
			}
			c.Orders(o)
		default:
		}
	}
}
