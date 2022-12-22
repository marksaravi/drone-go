package hardware

import (
	"log"
)

type gpio interface{}
type gpioreg interface{}
type i2c interface{}
type i2creg interface{}
type gpiooutput struct {
	pin gpio.PinOut
}

func (b *gpiooutput) SetHigh() {
	b.pin.Out(gpio.High)
}

func (b *gpiooutput) SetLow() {
	b.pin.Out(gpio.Low)
}

func NewGPIOOutput(pinName string) *gpiooutput {
	var pin gpio.PinOut = gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	dev := &gpiooutput{
		pin: pin,
	}
	dev.SetLow()
	return dev
}
