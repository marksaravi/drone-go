package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/esc"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	motor := flag.Int("motor", 0, "motor")
	flag.Parse()
	pca9685Configs := pca9685.ReadConfigs("./configs/hardware.json")
	powerBreakerGPIO := hardware.NewGPIOOutput(pca9685Configs.BreakerGPIO)
	powerBreaker := devices.NewPowerBreaker(powerBreakerGPIO)
	b, _ := i2creg.Open("/dev/i2c-1")
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}

	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:  i2cConn,
		MaxThrottle: pca9685Configs.MaxThrottle,
	})

	const maxThrottle float64 = 20
	const minThrottle float64 = 5
	const steps int = 10
	var dThrottle float64 = (maxThrottle - minThrottle) / float64(steps)
	motorsToChannelsMappings := make(map[int]int)
	for i := 0; i < 4; i++ {
		motorsToChannelsMappings[i] = pca9685Configs.MotorsMappings[i]
	}
	esc := esc.NewESC(pwmDev, motorsToChannelsMappings, powerBreaker, 50, false)
	var wg sync.WaitGroup
	esc.On()
	time.Sleep(5 * time.Second)
	throttle := float64(minThrottle)
	fmt.Println(*motor)
	motors := []float64{0, 0, 0, 0}
	for repeat := 0; repeat < 2; repeat++ {
		for step := 0; step < steps; step++ {
			// log.Println("motor: ", *motor, ", throttle:  ", throttle, "%")
			// for i:=0; i<len(throttles); i++ {
			// 	throttles[i]=throttle
			// }
			// esc.SetThrottles(throttles)
			motors[*motor] = throttle
			esc.SetThrottles(motors)
			time.Sleep(250 * time.Millisecond)
			throttle += dThrottle
		}
		dThrottle = -dThrottle
	}
	esc.Off()
	wg.Wait()
	log.Println("finished")
}
