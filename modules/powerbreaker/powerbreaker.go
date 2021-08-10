package powerbreaker

import (
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type PowerBreaker interface {
	Connect()
	Disconnect()
}

//powerBreaker is the safty power breaker
type powerBreaker struct {
	pin gpio.PinOut
}

//NewPowerBreaker creates new powerBreaker
func NewPowerBreaker(pinName string) *powerBreaker {
	var pin gpio.PinOut = gpioreg.ByName(pinName)
	pin.Out(gpio.Low)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return &powerBreaker{
		pin: pin,
	}
}

func (pb *powerBreaker) Connect() {
	pb.pin.Out(gpio.High)
	time.Sleep(4 * time.Second)
}

func (pb *powerBreaker) Disconnect() {
	pb.pin.Out(gpio.Low)
}
