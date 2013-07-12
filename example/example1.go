package main

import (
	"fmt"
	"github.com/Arrow/generals"
	"time"
)

const (
	width, height = 500, 500
)

type Frame struct{}
type Timing struct{}

type OrderTiming struct {
	o generals.Order
	t time.Duration
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
	rp := generals.RowRP(100, 250, 0.25, 2, 100)
	reg := generals.NewRegiment("20th Maine", rp, 500)

	orders := CallOrders([]OrderTiming{
		{generals.ForwardMarch, time.Second},
		{generals.Halt, 10 * time.Second},
		{generals.RightTurn, time.Second},
		{generals.RightTurn, time.Second},
		{generals.RightTurn, time.Second},
		{generals.ForwardMarch, time.Second},
		{generals.RightTurn, 5 * time.Second},
	})
	ctrlChan := SetupDisplay()
	FrameMsg := new(Frame)
	TimingMsg := new(Timing)
	tick := time.Tick(500 * time.Millisecond)
	timer := time.Tick(time.Second)

	for {

		ctrlChan <- reg

		ctrlChan <- FrameMsg
		<-tick

		select {
		case <-timer:
			//send timing message
			ctrlChan <- TimingMsg
		case o := <-orders:
			reg.Orders(o)
			fmt.Println("Order Received")
		default:
		}
	}
	time.Sleep(10 * time.Second)

}
