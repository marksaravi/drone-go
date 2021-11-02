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
	fmt.Println("Press ENTER to stop")
	hardware.InitHost()
	var pin gpio.PinOut = gpioreg.ByName("GPIO5")
	buzzer := piezobuzzer.NewBuzzer(pin)
	buzzer.Warning()
	fmt.Scanln()
	buzzer.Stop()
	time.Sleep(100 * time.Millisecond)
}
