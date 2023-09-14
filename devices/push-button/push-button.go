package pushbutton

import (
	"context"
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
	prevState gpio.Level
}

func NewPushButton(name string, pin gpioPin) *pushButton {
	return &pushButton{
		name:      name,
		pin:       pin,
		prevState: gpio.Low,
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
				state := b.pin.Read()
				if state == gpio.High && b.prevState == gpio.Low {
					ch <- true
				}
				b.prevState = state
			}
		}
		time.Sleep(time.Second / 100)
	}()
	return ch
}

func (b *pushButton) Name() string {
	return b.name
}
