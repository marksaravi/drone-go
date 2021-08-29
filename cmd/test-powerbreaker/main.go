package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers"
)

func main() {
	fmt.Println("Started")
	drivers.InitHost()
	powerbreaker := drivers.NewGPIOOutput("GPIO17")
	powerbreaker.SetHigh()
	time.Sleep(2 * time.Second)
	powerbreaker.SetLow()
}
