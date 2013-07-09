package generals

type Order int

const (
	Halt Order = iota
	ForwardMarch
	LeftWheel
	RightWheel
	LeftTurn
	RightTurn
)

type Orderer interface {
	Orders(Order)
}
