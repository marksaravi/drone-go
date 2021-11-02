package piezobuzzer

import (
	"context"
	"math"
	"time"

	"periph.io/x/periph/conn/gpio"
)

type SoundType struct {
	DevFrequency float64
	Steps        float64
	Duration     time.Duration
}

var Warning = SoundType{
	DevFrequency: 0,
	Steps:        25,
	Duration:     0,
}

var Siren = SoundType{
	DevFrequency: 200,
	Steps:        500,
	Duration:     0,
}

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

func (b *Buzzer) WaveGenerator(sound SoundType) {
	cx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	go func(ctx context.Context, buzzer *Buzzer) {
		buzzing := true
		const multiplier float64 = 1
		const baseFrequency float64 = 300 * multiplier
		var devFrequency float64 = sound.DevFrequency * multiplier
		const maxT float64 = 2
		const minT float64 = 1
		var dT = (maxT - minT) / sound.Steps //set to 500 for siren alarm
		var t float64 = minT
		on := true
		for buzzing {
			freq := baseFrequency + devFrequency*math.Exp(t)
			t += dT
			if t >= maxT {
				t = minT
				on = !on
			}
			select {
			case <-ctx.Done():
				buzzing = false
			default:
				if on {
					buzzer.out.Out(gpio.High)
				}
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
