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
	pulseMode bool
	pressed   bool
}

func NewPushButton(name string, pin gpioPin, pulseMode bool) *pushButton {
	return &pushButton{
		name: name,
		pin:  pin,
		pressed: false,
		pulseMode: pulseMode,
	}
}

func (b *pushButton) Name() string {
	return b.name
}

func (b *pushButton) PulseMode() bool {
	return b.pulseMode
}

func (b *pushButton) IsPressed() bool {
	pressed:= b.pin.Read() == gpio.Low
	if b.pulseMode {
		ispressed:=false
		if pressed && !b.pressed {
			ispressed=true
		}
		b.pressed=pressed
		return ispressed
	} else {
		return pressed
	}
}
