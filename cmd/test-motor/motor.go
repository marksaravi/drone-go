package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/devices/motors"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/pca9685"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func main() {
	drivers.InitHost()
	motor := flag.Int("motor", 0, "motor")
	flag.Parse()

	powerbreaker := drivers.NewGPIOOutput("GPIO17")
	b, _ := i2creg.Open("/dev/i2c-1")
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	escHandler, err := pca9685.NewPCA9685(pca9685.PCA9685Address, i2cConn, 15)
	if err != nil {
		log.Fatal(err)
	}
	motorsControl := motors.NewMotorsControl(escHandler, powerbreaker, map[int]int{0: 0, 1: 4, 2: 8, 3: 12})
	const maxThrottle float32 = 10
	const steps int = 10
	var dThrottle float32 = maxThrottle / float32(steps)
	var throttle float32 = 0
	motorsControl.On()
	time.Sleep(4 * time.Second)
	for repeat := 0; repeat < 4; repeat++ {
		for step := 0; step < steps; step++ {
			fmt.Println("motor: ", *motor, ", throttle:  ", throttle, "%")
			motorsControl.SetThrottle(*motor, throttle)
			time.Sleep(250 * time.Millisecond)
			throttle += dThrottle
		}
		dThrottle = -dThrottle
	}
	motorsControl.Off()
	fmt.Println("finished")
}
