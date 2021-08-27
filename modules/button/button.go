package button

import (
	"github.com/MarkSaravi/drone-go/hardware"
	"periph.io/x/periph/conn/gpio"
)

type button struct {
	pin gpio.PinIn
}

func NewButton(pinName string) *button {
	btn := hardware.NewButton(pinName)
	return &button{
		pin: btn,
	}
}

func (b *button) Read() bool {
	return b.pin.Read() == gpio.Low
}
