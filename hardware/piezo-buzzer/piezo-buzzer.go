package piezobuzzer

import (
	"context"
	"math"
	"sync"
	"time"

	"periph.io/x/periph/conn/gpio"
)

// CBGE

const (
	C      float64 = 16.35
	CSharp float64 = 17.32
	D      float64 = 18.35
	DSharp float64 = 19.45
	E      float64 = 20.60
	F      float64 = 21.83
	FSharp float64 = 23.12
	G      float64 = 24.50
	GSharp float64 = 25.96
	A      float64 = 27.50
	ASharp float64 = 29.14
	B      float64 = 30.87
)

type Note struct {
	Frequency float64
	Duration  time.Duration
	Octet     int
}

type Notes = []Note

type SoundWave struct {
	BaseFrequency float64
	DevFrequency  float64
	Steps         float64
	MinT          float64
	MaxT          float64
	Duration      time.Duration
}

var Warning = SoundWave{
	BaseFrequency: 300,
	DevFrequency:  0,
	Steps:         25,
	Duration:      0,
	MinT:          1,
	MaxT:          2,
}

var Siren = SoundWave{
	BaseFrequency: 300,
	DevFrequency:  200,
	Steps:         500,
	Duration:      0,
	MinT:          1,
	MaxT:          2,
}

var Connection1 = Notes{
	{Frequency: CSharp, Duration: time.Second / 6, Octet: 5},
	{Frequency: FSharp, Duration: time.Second / 6, Octet: 5},
}

type Buzzer struct {
	playing bool
	out     gpio.PinOut
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
}

func NewBuzzer(out gpio.PinOut) *Buzzer {
	var wg sync.WaitGroup
	buzzer := &Buzzer{
		playing: false,
		out:     out,
		cancel:  nil,
		wg:      &wg,
	}
	buzzer.out.Out(gpio.High)

	return buzzer
}

func (b *Buzzer) PlayNote(note Note) {
	start := time.Now()

	freq := note.Frequency * math.Pow(2, float64(note.Octet))
	period := time.Second / time.Duration(freq)
	for time.Since(start) < note.Duration {
		b.out.Out(gpio.High)
		onTime := time.Now()
		for time.Since(onTime) < 100*time.Microsecond {
		}
		b.out.Out(gpio.Low)
		for time.Since(onTime) < period {
		}
	}
}

func (b *Buzzer) PlayNotes(notes Notes) {
	for _, note := range notes {
		b.PlayNote(note)
	}
}

func (b *Buzzer) WaveGenerator(sound SoundWave) {
	if b.playing {
		return
	}
	cx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	go func(ctx context.Context, buzzer *Buzzer) {
		defer b.wg.Done()
		b.playing = true
		b.wg.Add(1)
		const multiplier float64 = 1
		var baseFrequency float64 = sound.BaseFrequency * multiplier
		var devFrequency float64 = sound.DevFrequency * multiplier
		var maxT float64 = sound.MaxT
		var minT float64 = sound.MinT
		var dT = (maxT - minT) / sound.Steps //set to 500 for siren alarm
		var t float64 = minT
		on := true
		for b.playing {
			freq := baseFrequency + devFrequency*math.Exp(t)
			t += dT
			if t >= maxT {
				t = minT
				on = !on
			}
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
			select {
			case <-ctx.Done():
				b.playing = false
			default:
			}
		}
		buzzer.out.Out(gpio.Low)
	}(cx, b)
}

func (b *Buzzer) Stop() {
	if b.cancel != nil {
		b.cancel()
		b.cancel = nil
	}
	b.wg.Wait()
}
