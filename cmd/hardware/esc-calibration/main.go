package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
)

func main() {
	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	breaker := powerbreaker.NewPowerBreaker()
	pwmDev.Start(float32(pca9685.Frequency))
	fmt.Println("setting max pulse width: ", pca9685.MaxPW)
	fmt.Println("turn on ESCs")
	pwmDev.SetPulseWidthAll(pca9685.MaxPW)
	time.Sleep(1 * time.Second)
	breaker.MotorsOn()
	time.Sleep(12 * time.Second)
	fmt.Println("setting min pulse width: ", pca9685.MinPW)
	pwmDev.SetPulseWidthAll(pca9685.MinPW)
	time.Sleep(12 * time.Second)
	fmt.Println("turn off ESCs")
	breaker.MotorsOff()
	time.Sleep(1 * time.Second)
	pwmDev.SetPulseWidthAll(0)
	pwmDev.Close()
}
