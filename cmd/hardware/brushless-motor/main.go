package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	channel := flag.Int("ch", 0, "ESC channel")
	flag.Parse()

	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}

	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection, 15, map[int]int{0: 0, 1: 4, 2: 8, 3: 12})
	if err != nil {
		fmt.Println(err)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	breaker := powerbreaker.NewPowerBreaker("GPIO17")

	pwmDev.Start()
	fmt.Println("channel: ", *channel)
	pwmDev.SetThrottle(*channel, 0)
	breaker.Connect()
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
	breaker.Disconnect()
	pwmDev.StopAll()
	pwmDev.Close()
	fmt.Println("finished")
}
