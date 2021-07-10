package powerbreaker

import (
	"github.com/MarkSaravi/drone-go/connectors/gpio"
)

//powerBreaker is the safty power breaker
type powerBreaker struct {
	gpio.GPIO
}

//NewPowerBreaker creates new powerBreaker
func NewPowerBreaker(pin gpio.GPIO) *powerBreaker {
	pin.SetAsOutput()
	pin.SetLow()
	return &powerBreaker{
		GPIO: pin,
	}
}
