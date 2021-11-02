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
	const dur = time.Second / 6
	const octet = 5
	buzzer.PlayNotes(piezobuzzer.Connection1)
}
