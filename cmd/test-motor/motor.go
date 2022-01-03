package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/esc"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()
	motor := flag.Int("motor", 0, "motor")
	flag.Parse()

	configs := config.ReadConfigs().FlightControl
	powerBreakerPin := configs.PowerBreaker
	powerBreakerGPIO := hardware.NewGPIOOutput(powerBreakerPin)
	powerBreaker := devices.NewPowerBreaker(powerBreakerGPIO)
	b, _ := i2creg.Open(configs.ESC.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Address, i2cConn, configs.ESC.MaxThrottle)

	const maxThrottle float64 = 10
	const steps int = 10
	var dThrottle float64 = maxThrottle / float64(steps)
	var throttle float64 = 0
	esc := esc.NewESC(pwmDev, powerBreaker, configs.ESC.PwmDeviceToESCMappings)
	var wg sync.WaitGroup
	esc.Start(&wg)
	esc.On()
	time.Sleep(3 * time.Second)
	throttles := map[uint8]float64{0: 0, 1: 0, 2: 0, 3: 0}
	for repeat := 0; repeat < 2; repeat++ {
		for step := 0; step < steps; step++ {
			log.Println("motor: ", *motor, ", throttle:  ", throttle, "%")
			throttles[uint8(*motor)] = throttle
			esc.SetThrottles(throttles)
			time.Sleep(250 * time.Millisecond)
			throttle += dThrottle
		}
		dThrottle = -dThrottle
	}
	esc.Off()
	esc.Close()
	wg.Wait()
	log.Println("finished")
}
