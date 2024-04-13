package hardware

import (
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

func NewPushButtonInput(gpioPin string) gpio.PinIn {
	pin := gpioreg.ByName(gpioPin)
	if pin == nil {
		log.Fatal("Failed to find ", gpioPin)
	}
	if err := pin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}
	return pin
}
