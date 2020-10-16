package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers/gpio"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
)

func main() {
	fmt.Println("Started")
	err := gpio.Open()
	defer gpio.Close()
	pin, err := gpio.NewPin(gpio.GPIO17)
	if err != nil {
		fmt.Println(err)
		return
	}
	breaker := powerbreaker.NewPowerBreaker(pin)
	defer breaker.SetLow()
	defer breaker.SetAsInput()
	time.Sleep(5 * time.Second)
}
