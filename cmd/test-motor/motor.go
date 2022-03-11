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
	"github.com/marksaravi/drone-go/models"
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

	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:      i2cConn,
		MaxThrottle:     15,
		ChannelMappings: configs.ESC.PwmDeviceToESCMappings,
	})

	const maxThrottle float64 = 10
	const minThrottle float64 = 5
	const steps int = 10
	var dThrottle float64 = (maxThrottle - minThrottle) / float64(steps)
	esc := esc.NewESC(pwmDev, powerBreaker, configs.ESC.UpdatePerSecond, false)
	var wg sync.WaitGroup
	esc.On()
	time.Sleep(5 * time.Second)
	throttles := models.Throttles{
		BaseThrottle: minThrottle,
		Throttles:    map[int]float64{0: 0, 1: 0, 2: 0, 3: 0},
	}
	for repeat := 0; repeat < 2; repeat++ {
		for step := 0; step < steps; step++ {
			log.Println("motor: ", *motor, ", throttle:  ", throttles.BaseThrottle, "%")
			throttles.Throttles[*motor] = throttles.BaseThrottle
			esc.SetThrottles(throttles)
			time.Sleep(250 * time.Millisecond)
			throttles.BaseThrottle += dThrottle
		}
		dThrottle = -dThrottle
	}
	esc.Off()
	wg.Wait()
	log.Println("finished")
}
