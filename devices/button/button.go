package button

import (
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const ON_LEVEL = gpio.Low
const OFF_LEVEL = gpio.High

type Button struct {
	name string
	pin  gpio.PinIn
	prev gpio.Level
}

func NewButton(gpioName string) *Button {
	pin := gpioreg.ByName(gpioName)
	if pin == nil {
		log.Fatal("Failed to find ", gpioName)
	}
	return &Button{
		name: gpioName,
		pin:  pin,
		prev: OFF_LEVEL,
	}
}

func (b *Button) Read() (level bool, pressed bool) {
	l := b.pin.Read()
	pressed = b.prev == OFF_LEVEL && l == ON_LEVEL
	b.prev = l
	level = l == ON_LEVEL
	return level, pressed
}

func (b *Button) Name() string {
	return b.name
}
