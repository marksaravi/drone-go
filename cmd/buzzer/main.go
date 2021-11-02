package main

import (
	"fmt"

	"github.com/marksaravi/drone-go/hardware"
	piezobuzzer "github.com/marksaravi/drone-go/hardware/piezo-buzzer"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

func main() {
	fmt.Println("Press ENTER to change sounds")
	hardware.InitHost()
	var pin gpio.PinOut = gpioreg.ByName("GPIO5")
	buzzer := piezobuzzer.NewBuzzer(pin)
	buzzer.WaveGenerator(piezobuzzer.Warning)
	fmt.Scanln()
	buzzer.Stop()
	buzzer.WaveGenerator(piezobuzzer.Siren)
	fmt.Scanln()
	buzzer.Stop()
}
