package pushbutton

import (
	"context"
	"time"
)

const EDGE_TIMEOUT = time.Millisecond * 10

type gpioPin interface {
	WaitForEdge(timeout time.Duration) bool
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
