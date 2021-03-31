package powerbreaker

import (
	"github.com/MarkSaravi/drone-go/connectors/gpio"
)

//PowerBreaker is the safty power breaker
type PowerBreaker struct {
	gpio.GPIO
}

//NewPowerBreaker creates new PowerBreaker
func NewPowerBreaker(pin gpio.GPIO) *PowerBreaker {
	pin.SetAsOutput()
	pin.SetLow()
	return &PowerBreaker{
		GPIO: pin,
	}
}
