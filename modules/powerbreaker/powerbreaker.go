package powerbreaker

import (
	"github.com/MarkSaravi/drone-go/connectors/gpio"
	"github.com/MarkSaravi/drone-go/types"
)

//powerBreaker is the safty power breaker
type powerBreaker struct {
	pin types.GPIO
}

//NewPowerBreaker creates new powerBreaker
func NewPowerBreaker() *powerBreaker {
	gpio.Open()
	pin, _ := gpio.NewPin(gpio.GPIO17)
	pin.SetAsInput()
	return &powerBreaker{
		pin: pin,
	}
}

func (pb *powerBreaker) MotorsOn() {
	pb.pin.SetAsOutput()
	pb.pin.SetHigh()
}

func (pb *powerBreaker) MotorsOff() {
	pb.pin.SetLow()
	pb.pin.SetAsInput()
	gpio.Close()
}
