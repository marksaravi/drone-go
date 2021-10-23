package hardware

import (
	"fmt"
	"log"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type gpioswitch struct {
	pin gpio.PinIn
}

func (b *gpioswitch) Read() bool {
	return b.pin.Read() == gpio.Low
}

func NewGPIOSwitch(pinName string) *gpioswitch {
	fmt.Println("GPIO: ", pinName)
	var pin gpio.PinIn = gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	pin.In(gpio.Float, gpio.NoEdge)
	return &gpioswitch{
		pin: pin,
	}
}
