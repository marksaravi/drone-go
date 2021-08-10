package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/modules/motors"
)

func main() {
	appConfig := config.ReadConfigs()

	_, esc, _, powerbreaker := hardware.InitDroneHardware(appConfig)
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
