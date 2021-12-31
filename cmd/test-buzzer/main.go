package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/marksaravi/drone-go/hardware"
	piezobuzzer "github.com/marksaravi/drone-go/hardware/piezo-buzzer"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

func main() {
	hardware.InitHost()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	var pin gpio.PinOut = gpioreg.ByName("GPIO5")
	buzzer := piezobuzzer.NewBuzzer(pin)
	fmt.Println("Warning sound, press ENTER to next sound")
	buzzer.WaveGenerator(ctx, &wg, piezobuzzer.WarningSound)
	fmt.Scanln()
	cancel()
	wg.Wait()
	ctx, cancel = context.WithCancel(context.Background())
	fmt.Println("Siren sound, press ENTER to stop")
	buzzer.WaveGenerator(ctx, &wg, piezobuzzer.SirenSound)
	fmt.Scanln()
	buzzer.Stop()
	buzzer.PlaySound(piezobuzzer.ConnectedSound)
	fmt.Scanln()
	buzzer.PlaySound(piezobuzzer.DisconnectedSound)
	cancel()
}
