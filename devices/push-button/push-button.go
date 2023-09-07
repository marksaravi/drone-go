package pushbutton

import (
	"context"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const EDGE_TIMEOUT = time.Millisecond * 10

type Button struct {
	gpioPin string
	name    string
	pin     gpio.PinIn
}

func NewPushButton(gpioPin string, name string, pullUp bool) *Button {
	pin := gpioreg.ByName(gpioPin)
	if pin == nil {
		log.Fatal("Failed to find ", gpioPin)
	}
	if err := pin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}

	return &Button{
		gpioPin: gpioPin,
		name:    name,
		pin:     pin,
	}
}

func (b *Button) Start(ctx context.Context) <-chan bool {
	ch := make(chan bool, 1)

	go func() {
		for ch != nil {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				if b.IsPushed() {
					ch <- true
				}
			}
		}
	}()
	return ch
}

func (b *Button) IsPushed() bool {
	return b.pin.WaitForEdge(EDGE_TIMEOUT)
}
func (b *Button) Name() string {
	return b.name
}

func (b *Button) GPIO() string {
	return b.gpioPin
}
