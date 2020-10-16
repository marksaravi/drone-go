package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers/gpio"
)

func main() {
	pn := flag.Int("pin", 2, "Pin")
	flag.Parse()
	gpiopin := gpio.GPIO02
	switch *pn {
	case 4:
		gpiopin = gpio.GPIO04
	case 17:
		gpiopin = gpio.GPIO17
	default:
		gpiopin = 2
	}
	err := gpio.Open()
	defer gpio.Close()
	fmt.Println("GPIO is opened successfully")
	pin, err := gpio.NewPin(gpiopin)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Set output")
	pin.SetAsOutput()
	defer pin.SetAsInput()
	fmt.Println("Set High")
	pin.SetHigh()
	time.Sleep(5 * time.Second)
	fmt.Println("Set Low")
	pin.SetLow()
}
