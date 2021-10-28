package hardware

import (
	"log"

	"github.com/marksaravi/drone-go/config"
	"periph.io/x/periph/host"
)

func InitHost() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}

func NewPowerBreaker() interface {
	SetLow()
	SetHigh()
} {
	flightControlConfig := config.ReadFlightControlConfig()
	powerbreaker := NewGPIOOutput(flightControlConfig.Configs.PowerBreaker)
	return powerbreaker
}