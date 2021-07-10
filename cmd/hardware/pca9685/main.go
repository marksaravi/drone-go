package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
)

func main() {
	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	pca9685, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = pca9685.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pca9685.Close()
	pca9685.SetPulseWidth(0, 0.002)
	pca9685.SetPulseWidth(1, 0.002)
	pca9685.SetPulseWidth(2, 0.002)
	pca9685.SetPulseWidth(3, 0.002)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Channel 0 PWM is set")
	}
	time.Sleep(5 * time.Second)
	pca9685.StopAll()
}
