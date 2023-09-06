package pushbutton

import (
	"context"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const TIMEOUT = time.Millisecond * 200

type Button struct {
	gpioPin  string
	name     string
	pin      gpio.PinIn
	lastPush time.Time
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
		b.lastPush = time.Now()
		for ch != nil {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				for b.pin.WaitForEdge(time.Millisecond * 10) {
					if time.Since(b.lastPush) >= TIMEOUT {
						ch <- true
						b.lastPush = time.Now()
					}
				}
			}
		}
	}()
	return ch
}

func (b *Button) Name() string {
	return b.name
}

func (b *Button) GPIO() string {
	return b.gpioPin
}
