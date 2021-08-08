package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/motors"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
	"periph.io/x/periph/host"
)

func initHardware(config types.EscConfig) (types.ESC, types.PowerBreaker) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	i2cConnection, err := i2c.Open(config.Device)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection, 15, config.Motors)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	powerbreaker := powerbreaker.NewPowerBreaker(config.PowerBrokerGPIO)
	pwmDev.Start()
	pwmDev.StopAll()
	return pwmDev, powerbreaker
}

func main() {
	appConfig := utils.ReadConfigs()

	esc, powerbreaker := initHardware(appConfig.Flight.Esc)
	motor := flag.Int("motor", 0, "motor")
	flag.Parse()

	motorsControl := motors.NewMotorsControl(esc, powerbreaker)
	time.Sleep(4 * time.Second)
	const maxThrottle float32 = 10
	const steps int = 10
	var dThrottle float32 = maxThrottle / float32(steps)
	var throttle float32 = 0
	motorsControl.On()
	throttles := map[int]float32{0: 0, 1: 0, 2: 0, 3: 0}
	for repeat := 0; repeat < 4; repeat++ {
		for step := 0; step < steps; step++ {
			fmt.Println("motor: ", *motor, ", throttle:  ", throttle, "%")
			throttles[*motor] = throttle
			motorsControl.SetThrottles(throttles)
			time.Sleep(250 * time.Millisecond)
			throttle += dThrottle
		}
		dThrottle = -dThrottle
	}
	motorsControl.Off()
	fmt.Println("finished")
}
