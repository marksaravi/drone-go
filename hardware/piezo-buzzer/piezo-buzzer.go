package piezobuzzer

import (
	"math"
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

var WarningSound = SoundWave{
	BaseFrequency: 300,
	DevFrequency:  0,
	Steps:         25,
	Duration:      0,
	MinT:          1,
	MaxT:          2,
}

var SirenSound = SoundWave{
	BaseFrequency: 300,
	DevFrequency:  200,
	Steps:         500,
	Duration:      0,
	MinT:          1,
	MaxT:          2,
}

var ConnectedSound = Notes{
	{Frequency: C, Duration: time.Second / 6, Octet: 4},
	{Frequency: G, Duration: time.Second / 6, Octet: 4},
}

var DisconnectedSound = Notes{
	{Frequency: G, Duration: time.Second / 6, Octet: 4},
	{Frequency: C, Duration: time.Second / 6, Octet: 4},
}

type Buzzer struct {
	playing int
	out     gpio.PinOut
}

func NewBuzzer(out gpio.PinOut) *Buzzer {
	buzzer := &Buzzer{
		playing: 0,
		out:     out,
	}
	buzzer.out.Out(gpio.High)

	return buzzer
}

func (b *Buzzer) playNote(note Note) {
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

func (b *Buzzer) PlaySound(notes Notes) {
	b.Stop()
	go func() {
		b.playNotes(notes)
	}()
}

func (b *Buzzer) playNotes(notes Notes) {
	for _, note := range notes {
		b.playNote(note)
	}
}

func (b *Buzzer) WaveGenerator(sound SoundWave) {
	b.playing++
	id := b.playing
	go func() {
		const multiplier float64 = 1
		var baseFrequency float64 = sound.BaseFrequency * multiplier
		var devFrequency float64 = sound.DevFrequency * multiplier
		var maxT float64 = sound.MaxT
		var minT float64 = sound.MinT
		var dT = (maxT - minT) / sound.Steps //set to 500 for siren alarm
		var t float64 = minT
		on := true
		for b.playing == id {
			freq := baseFrequency + devFrequency*math.Exp(t)
			t += dT
			if t >= maxT {
				t = minT
				on = !on
			}
			if on {
				b.out.Out(gpio.High)
			}
			period := time.Second / time.Duration(freq)
			onTime := time.Now()
			for time.Since(onTime) < 100*time.Microsecond {
			}
			b.out.Out(gpio.Low)
			for time.Since(onTime) < period {
			}
		}
		b.out.Out(gpio.Low)
	}()
}

func (b *Buzzer) Stop() {
	b.playing++
}
