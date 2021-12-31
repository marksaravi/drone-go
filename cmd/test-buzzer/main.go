package main

import (
	"fmt"

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
	buzzer.WaveGenerator(piezobuzzer.WarningSound)
	fmt.Scanln()
	buzzer.Stop()
	fmt.Println("Siren sound, press ENTER to stop")
	buzzer.WaveGenerator(piezobuzzer.SirenSound)
	fmt.Scanln()
	buzzer.Stop()
	buzzer.PlaySound(piezobuzzer.ConnectedSound)
	fmt.Scanln()
	buzzer.PlaySound(piezobuzzer.DisconnectedSound)
}
