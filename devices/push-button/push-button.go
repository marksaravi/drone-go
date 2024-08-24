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
	name         string
	index        int
	isPushButton bool
	pin          gpioPin
	wasPressed   bool
	pressLock    bool
}

func NewPushButton(name string, index int, isPushButton bool, pin gpioPin) *pushButton {
	return &pushButton{
		name: name,
		index: index,
		isPushButton: isPushButton,
		pin:  pin,
		wasPressed: false,
		pressLock: false,
	}
}

func (b *pushButton) Name() string {
	return b.name
}

func (b *pushButton) Index() int {
	return b.index
}

func (b *pushButton) IsPushButton() bool {
	return b.isPushButton
}


func (b *pushButton) IsPressed() bool {
	return b.pin.Read() == gpio.Low
}

func (b *pushButton) Update() {
	if !b.wasPressed && !b.pressLock && b.IsPressed() {
		b.wasPressed = true
		b.pressLock = true
	}
}

func (b *pushButton) IsPushed() bool {
	if !b.IsPressed() && b.pressLock {
		b.pressLock = false
	}	
	if b.wasPressed {
		b.wasPressed = false
		return true
	}
	return false
}