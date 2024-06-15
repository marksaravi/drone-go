package main

import (
	"context"
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
	b, _ := i2creg.Open(pca9685Configs.I2CPort)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}

	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:      i2cConn,
		MaxSafeThrottle: pca9685Configs.MaxSafeThrottle,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Scanln()
		cancel()
	}()

	const maxThrottle float64 = 80
	const minThrottle float64 = 5
	const steps int = 40
	var dThrottle float64 = (maxThrottle - minThrottle) / float64(steps)
	esc := esc.NewESC(pwmDev, pca9685Configs.MotorsMappings, powerBreaker, 50, false)
	var wg sync.WaitGroup
	esc.On()
	time.Sleep(5 * time.Second)
	throttle := float64(minThrottle)
	fmt.Println(*motor)
	motors := []float64{0, 0, 0, 0}
	running := true
	const runninTimeMilliSecond int = 600 * 1000
	const runninTimeDur int = 250
	const numOfInnerLoop int = runninTimeMilliSecond / runninTimeDur
	for repeat := 0; repeat < 2 && running; repeat++ {

		for step := 0; step < steps && running; step++ {
			select {
			case _, running = <-ctx.Done():
			default:
			}
			motors[*motor] = throttle
			esc.SetThrottles(motors)
			time.Sleep(100 * time.Millisecond)
			throttle += dThrottle
		}
		if throttle == maxThrottle {
			for i := 0; i < numOfInnerLoop && running; i++ {
				select {
				case _, running = <-ctx.Done():
				default:
				}
				motors[*motor] = throttle
				esc.SetThrottles(motors)
				time.Sleep(time.Duration(runninTimeDur) * time.Millisecond)
			}
		}
		dThrottle = -dThrottle

	}
	esc.Off()
	wg.Wait()
	log.Println("finished")
}
