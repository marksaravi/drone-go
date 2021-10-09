package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/devices/motors"
	"github.com/marksaravi/drone-go/hardware"
)

func main() {
	hardware.InitHost()
	motor := flag.Int("motor", 0, "motor")
	flag.Parse()

	const maxThrottle float32 = 10
	const steps int = 10
	var dThrottle float32 = maxThrottle / float32(steps)
	var throttle float32 = 0
	esc := motors.NewESC()
	esc.On()
	time.Sleep(3 * time.Second)
	throttles := map[uint8]float32{0: 0, 1: 0, 2: 0, 3: 0}
	for repeat := 0; repeat < 2; repeat++ {
		for step := 0; step < steps; step++ {
			fmt.Println("motor: ", *motor, ", throttle:  ", throttle, "%")
			throttles[uint8(*motor)] = throttle
			s := time.Now()
			esc.SetThrottles(throttles)
			fmt.Println(time.Since(s))
			time.Sleep(250 * time.Millisecond)
			throttle += dThrottle
		}
		dThrottle = -dThrottle
	}
	esc.Off()
	fmt.Println("finished")
}
