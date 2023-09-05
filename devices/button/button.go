package button

import (
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const ON_LEVEL = gpio.Low
const OFF_LEVEL = gpio.High

type Button struct {
	gpio string
	name string
	pin  gpio.PinIn
	prev gpio.Level
}

func NewButton(gpio string, name string) *Button {
	pin := gpioreg.ByName(gpio)
	if pin == nil {
		log.Fatal("Failed to find ", gpio)
	}
	return &Button{
		gpio: gpio,
		name: name,
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

func (b *Button) GPIO() string {
	return b.gpio
}
