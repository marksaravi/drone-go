package pushbutton

import (
	"context"
	"time"

	"periph.io/x/conn/v3/gpio"
)

const EDGE_TIMEOUT = time.Millisecond * 10

type gpioPin interface {
	Read() gpio.Level
	WaitForEdge(duration time.Duration) bool
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

func (b *pushButton) Start(ctx context.Context) <-chan bool {
	ch := make(chan bool, 1)

	go func() {
		for ch != nil {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				if b.pin.WaitForEdge(time.Second / 20) {
					ch <- true
				}
			}
		}
	}()
	return ch
}

func (b *pushButton) Name() string {
	return b.name
}
