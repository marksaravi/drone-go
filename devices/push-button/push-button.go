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
	name    string
	pin     gpioPin
	hold    bool
	pressed bool
}

func NewPushButton(name string, pin gpioPin, hold bool) *pushButton {
	return &pushButton{
		name: name,
		pin:  pin,
		pressed: false,
		hold: hold,
	}
}

func (b *pushButton) Name() string {
	return b.name
}

func (b *pushButton) Hold() bool {
	return b.hold
}

func (b *pushButton) IsPressed() bool {
	pressed:= b.pin.Read() == gpio.Low
	if b.hold {
		return pressed
	} else {
		ispressed:=false
		if pressed && !b.pressed {
			ispressed=true
		}
		b.pressed=pressed
		return ispressed
	}
}
