package pushbutton

import (
	"context"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const EDGE_TIMEOUT = time.Millisecond * 10

type pushButton struct {
	gpioPin string
	name    string
	pin     gpio.PinIn
}

func NewPushButton(name string, gpioPin string) *pushButton {
	pin := gpioreg.ByName(gpioPin)
	if pin == nil {
		log.Fatal("Failed to find ", name)
	}
	if err := pin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}

	return &pushButton{
		gpioPin: gpioPin,
		name:    name,
		pin:     pin,
	}
}

func (b *pushButton) Start(ctx context.Context) <-chan bool {
	ch := make(chan bool, 1)

	go func() {
		for ch != nil {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				if b.WaitForPush() {
					ch <- true
				}
			}
		}
	}()
	return ch
}

func (b *pushButton) WaitForPush() bool {
	return b.pin.WaitForEdge(EDGE_TIMEOUT)
}

func (b *pushButton) Name() string {
	return b.name
}

func (b *pushButton) GPIO() string {
	return b.gpioPin
}
