package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devices/pca9685"
	"github.com/MarkSaravi/drone-go/drivers/i2c"
	"github.com/MarkSaravi/drone-go/modules/esc"
	"github.com/MarkSaravi/drone-go/drivers/gpio"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
)

func main() {
	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	pca9685, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection)
	esc := esc.NewESC(pca9685)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	err = gpio.Open()
	defer gpio.Close()
	pin, err := gpio.NewPin(gpio.GPIO17)
	if err != nil {
		fmt.Println(err)
		return
	}
	breaker := powerbreaker.NewPowerBreaker(pin)
	defer breaker.SetLow()
	defer breaker.SetAsInput()

	esc.Start(float32(50))
	esc.StopAll()
	defer esc.Close()
	fmt.Println("setting pulse width to 0.002")
	fmt.Println("turn on ESCs")
	esc.SetPulseWidthAll(0.002)
	breaker.SetHigh()
	time.Sleep(10 * time.Second)
	fmt.Println("setting pulse width to 0.001")
	esc.SetPulseWidthAll(0.001)
	time.Sleep(10 * time.Second)
	fmt.Println("turn off ESCs")
	breaker.SetLow()
	time.Sleep(1 * time.Second)
	fmt.Println("setting pulse with to 0")
	esc.StopAll()

}
