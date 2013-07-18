package generals

type Order int

const (
	Halt Order = iota
	ForwardMarch
	LeftWheel
	RightWheel
	LeftTurn
	RightTurn
	Reform
	FormRow
	FormTwoRow
	FormThreeRow
	FormFourRow
	FormCol
	FormTwoCol
	FormThreeCol
	FormFourCol
	FormFiveCol
	Quit
	PrintForm
)

type Orderer interface {
	Orders(Order)
}
