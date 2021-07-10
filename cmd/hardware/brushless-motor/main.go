package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/connectors/gpio"
	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
)

func main() {
	channel := flag.Int("ch", 0, "ESC channel")
	flag.Parse()

	maxPower := 25

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
	defer pwmDev.SetPulseWidthAll(0)
	defer pwmDev.Close()

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

	pwmDev.Start(pca9685.Frequency)
	fmt.Println("channel: ", *channel, ", PW:  ", pca9685.MinPW)
	pwmDev.SetPulseWidth(*channel, pca9685.MinPW)
	breaker.SetHigh()
	time.Sleep(4 * time.Second)
	power := 1
	inc := 1
	for power != 0 {
		pw := float32(pca9685.MaxPW-pca9685.MinPW)*float32(power)/100 + float32(pca9685.MinPW)
		fmt.Println("channel: ", *channel, ", PW:  ", pw)
		pwmDev.SetPulseWidth(*channel, pw)
		time.Sleep(250 * time.Millisecond)
		power += inc
		if power == maxPower {
			inc = -1
		}
	}
	breaker.SetLow()
	fmt.Println("finished")
}
