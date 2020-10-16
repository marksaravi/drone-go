package powerbreaker

import (
	"github.com/MarkSaravi/drone-go/drivers/gpio"
)

//PowerBreaker is the safty power breaker
type PowerBreaker struct {
	gpio.GPIO
}

//NewPowerBreaker creates new PowerBreaker
func NewPowerBreaker(pin gpio.GPIO) *PowerBreaker {
	pin.SetAsOutput()
	pin.SetHigh()
	return &PowerBreaker{
		GPIO: pin,
	}
}
