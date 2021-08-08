package powerbreaker

import (
	"log"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

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

func (pb *powerBreaker) MotorsOn() {
	pb.pin.Out(gpio.High)
}

func (pb *powerBreaker) MotorsOff() {
	pb.pin.Out(gpio.Low)
}
