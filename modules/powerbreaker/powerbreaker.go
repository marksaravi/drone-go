package powerbreaker

import (
	"github.com/MarkSaravi/drone-go/connectors/gpio"
)

//powerBreaker is the safty power breaker
type powerBreaker struct {
	pin *gpio.Pin
}

//NewPowerBreaker creates new powerBreaker
func NewPowerBreaker() *powerBreaker {
	gpio.Open()
	pin, _ := gpio.NewPin(gpio.GPIO17)
	pin.SetAsOutput()
	pin.SetLow()
	return &powerBreaker{
		pin: pin,
	}
}

func (pb *powerBreaker) MototsOn() {
	pb.pin.SetAsOutput()
	pb.pin.SetHigh()
}

func (pb *powerBreaker) MototsOff() {
	pb.pin.SetLow()
	pb.pin.SetAsInput()
	gpio.Close()
}
