package piezobuzzer

import (
	"context"
	"math"
	"time"

	"periph.io/x/periph/conn/gpio"
)

type Buzzer struct {
	out    gpio.PinOut
	cancel context.CancelFunc
}

func NewBuzzer(out gpio.PinOut) *Buzzer {
	buzzer := &Buzzer{
		out:    out,
		cancel: nil,
	}
	buzzer.out.Out(gpio.High)

	return buzzer
}

func (b *Buzzer) Buzz() {
	cx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	go func(ctx context.Context, buzzer *Buzzer) {
		buzzing := true
		const baseFrequency float64 = 300
		const devFrequency float64 = 200
		var t float64 = 1
		for buzzing {
			freq := baseFrequency + devFrequency*math.Exp(t)
			t += 0.005
			if t >= 2 {
				t = 1
			}
			select {
			case <-ctx.Done():
				buzzing = false
			default:
				buzzer.out.Out(gpio.High)
				period := time.Second / time.Duration(freq)
				onTime := time.Now()
				for time.Since(onTime) < 100*time.Microsecond {

				}
				buzzer.out.Out(gpio.Low)
				for time.Since(onTime) < period {

				}
			}
		}
		buzzer.out.Out(gpio.Low)
	}(cx, b)
}

func (b *Buzzer) Stop() {
	if b.cancel != nil {
		b.cancel()
	}
}
