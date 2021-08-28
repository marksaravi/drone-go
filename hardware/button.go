package hardware

import (
	"log"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type button struct {
	pin   gpio.PinIn
	value bool
}

func (b *button) Read() bool {
	return b.pin.Read() == gpio.Low
}

func (b *button) Value() bool {
	return b.value
}

func NewButton(pinName string) *button {
	var pin gpio.PinIn = gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	pin.In(gpio.Float, gpio.NoEdge)
	return &button{
		pin:   pin,
		value: false,
	}
}
