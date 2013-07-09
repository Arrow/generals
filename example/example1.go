package main

import (
	"fmt"
	"time"
	"github.com/Arrow/generals"
)

const (
	width, height    = 500, 500
)

type Frame struct { }
type Timing struct { }

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
	/*rpc := new(generals.RegPlacement)
	//rp.Rows = 20
	rpc.RowPath = []generals.Point{
		generals.Point{100, 250},
		generals.Point{95, 250},
		generals.Point{90, 250},
		generals.Point{85, 250},
		generals.Point{80, 250},
		generals.Point{75, 250},
		generals.Point{70, 250},
		generals.Point{65, 250},
		generals.Point{60, 250},
		generals.Point{55, 250},
		generals.Point{50, 250},
		generals.Point{45, 250},
		generals.Point{40, 250},
		generals.Point{35, 250},
		generals.Point{30, 250},
		generals.Point{25, 250},
		generals.Point{20, 250},
		generals.Point{15, 250},
		generals.Point{10, 250},
		generals.Point{5, 250},
	}
	//rp.Columns = 5
	rpc.ColPath = []generals.Point{
		generals.Point{100, 250},
		generals.Point{100, 245},
		generals.Point{100, 240},
		generals.Point{100, 235},
		generals.Point{100, 230},
	}
	*/
	rp := generals.RowRP(100, 250, 0.25, 2, 100)
	
	//fmt.Println(rpc)
	//fmt.Println(rp)
	
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
