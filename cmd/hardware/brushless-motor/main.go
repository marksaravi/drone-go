package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
)

func main() {
	channel := flag.Int("ch", 0, "ESC channel")
	flag.Parse()

	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}

	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection, 15, map[int]int{0: 5, 1: 6, 2: 7, 3: 13})
	if err != nil {
		fmt.Println(err)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	breaker := powerbreaker.NewPowerBreaker()

	pwmDev.Start()
	fmt.Println("channel: ", *channel)
	pwmDev.SetThrottle(*channel, 0)
	breaker.MotorsOn()
	time.Sleep(4 * time.Second)
	const maxThrottle float32 = 10
	const steps int = 10
	var dThrottle float32 = maxThrottle / float32(steps)
	var throttle float32 = 0
	for repeat := 0; repeat < 4; repeat++ {
		for step := 0; step < steps; step++ {
			fmt.Println("channel: ", *channel, ", throttle:  ", throttle, "%")
			pwmDev.SetThrottle(*channel, throttle)
			time.Sleep(250 * time.Millisecond)
			throttle += dThrottle
		}
		dThrottle = -dThrottle
	}
	breaker.MotorsOff()
	pwmDev.StopAll()
	pwmDev.Close()
	fmt.Println("finished")
}
