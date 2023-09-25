package pushbutton

import (
	"time"

	"periph.io/x/conn/v3/gpio"
)

const EDGE_TIMEOUT = time.Millisecond * 10

type gpioPin interface {
	Read() gpio.Level
}

type pushButton struct {
	name string
	pin  gpioPin
}

func NewPushButton(name string, pin gpioPin) *pushButton {
	return &pushButton{
		name: name,
		pin:  pin,
	}
}

func (b *pushButton) Name() string {
	return b.name
}

func (b *pushButton) Read() bool {
	return b.pin.Read() == gpio.High
}
