package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/hardware"
)

func main() {
	fmt.Println("Started")
	hardware.HostInitialize()
	powerbreakerGPIO := config.ReadConfigs().FlightControl.PowerBreaker
	gpiopin := hardware.NewGPIOOutput(powerbreakerGPIO)
	powerbreaker := devices.NewPowerBreaker(gpiopin)
	powerbreaker.Connect()
	time.Sleep(3 * time.Second)
	powerbreaker.Disconnect()
}
