package types

//GPIO is interface to GPIO pin
type GPIO interface {
	SetAsInput()
	SetAsOutput()
	SetHigh()
	SetLow()
}

type Throttle struct {
	Motor int
	Value float32
}
