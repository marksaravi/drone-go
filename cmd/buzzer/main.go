package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	piezobuzzer "github.com/marksaravi/drone-go/hardware/piezo-buzzer"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

func main() {
	hardware.InitHost()
	var pin gpio.PinOut = gpioreg.ByName("GPIO5")
	buzzer := piezobuzzer.NewBuzzer(pin)
	fmt.Println("Warning sound, press ENTER to next sound")
	// buzzer.WaveGenerator(piezobuzzer.Warning)
	// fmt.Scanln()
	// buzzer.Stop()
	// fmt.Println("Siren sound, press ENTER to stop")
	// buzzer.WaveGenerator(piezobuzzer.Siren)
	// fmt.Scanln()
	// buzzer.Stop()
	const dur = time.Second / 2
	const octet = 3
	notes := piezobuzzer.Notes{
		{Frequency: piezobuzzer.C, Duration: dur, Octet: octet},
		{Frequency: piezobuzzer.D, Duration: dur, Octet: octet},
		{Frequency: piezobuzzer.E, Duration: dur, Octet: octet},
		{Frequency: piezobuzzer.F, Duration: dur, Octet: octet},
		{Frequency: piezobuzzer.G, Duration: dur, Octet: octet},
		{Frequency: piezobuzzer.A, Duration: dur, Octet: octet},
		{Frequency: piezobuzzer.B, Duration: dur, Octet: octet},
	}
	for i := octet; i < octet+4; i++ {
		for n := 0; n < len(notes); n++ {
			notes[n].Octet = i
			buzzer.PlayNote(notes[n])
		}
	}
	buzzer.Stop()
}
