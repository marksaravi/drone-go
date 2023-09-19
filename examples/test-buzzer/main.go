package main

import (
	"fmt"

	"github.com/marksaravi/drone-go/hardware"
	piezobuzzer "github.com/marksaravi/drone-go/hardware/piezo-buzzer"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

func main() {
	hardware.HostInitialize()
	var pin gpio.PinOut = gpioreg.ByName("GPIO5")
	buzzer := piezobuzzer.NewBuzzer(pin)
	fmt.Println("Warning sound, press ENTER to next sound")
	buzzer.WaveGenerator(piezobuzzer.WarningSound)
	fmt.Scanln()
	fmt.Println("Siren sound, press ENTER to next sound")
	buzzer.WaveGenerator(piezobuzzer.SirenSound)
	fmt.Scanln()
	fmt.Println("Connected sound")
	buzzer.PlaySound(piezobuzzer.ConnectedSound)
	fmt.Scanln()
	fmt.Println("Disconnected sound")
	buzzer.PlaySound(piezobuzzer.DisconnectedSound)
	fmt.Scanln()
}
