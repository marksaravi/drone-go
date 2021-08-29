package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	drivers.InitHost()
	motor := flag.Int("motor", 0, "motor")
	flag.Parse()

	const maxThrottle float32 = 10
	const steps int = 10
	var dThrottle float32 = maxThrottle / float32(steps)
	var throttle float32 = 0
	esc := utils.NewESC()
	esc.On()
	time.Sleep(4 * time.Second)
	for repeat := 0; repeat < 4; repeat++ {
		for step := 0; step < steps; step++ {
			fmt.Println("motor: ", *motor, ", throttle:  ", throttle, "%")
			esc.SetThrottle(*motor, throttle)
			time.Sleep(250 * time.Millisecond)
			throttle += dThrottle
		}
		dThrottle = -dThrottle
	}
	esc.Off()
	fmt.Println("finished")
}
