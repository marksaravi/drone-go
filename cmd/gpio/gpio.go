package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers/gpio"
)

func main() {
	fmt.Println("Testing my GPIO.")
	err := gpio.Open()
	defer gpio.Close()
	fmt.Println("GPIO is opened successfully")
	pin, err := gpio.NewPin(gpio.GPIO04)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Set output")
	pin.SetOutput()
	defer pin.SetInput()
	fmt.Println("Set High")
	pin.SetHigh()
	time.Sleep(5 * time.Second)
	fmt.Println("Set Low")
	pin.SetLow()
}
