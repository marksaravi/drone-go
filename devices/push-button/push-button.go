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
	name      string
	pin       gpioPin
	wasPressed   bool
}

func NewPushButton(name string, pin gpioPin) *pushButton {
	return &pushButton{
		name: name,
		pin:  pin,
		wasPressed: false,
	}
}

func (b *pushButton) Name() string {
	return b.name
}

func (b *pushButton) IsPressed() bool {
	return b.pin.Read() == gpio.Low
}

func (b *pushButton) Update() {
	if !b.wasPressed && b.IsPressed() {
		b.wasPressed = true
	}
}

func (b *pushButton) IsPushed() bool {
	if b.wasPressed {
		b.wasPressed = false
		return true
	}
	return false
}