package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/hardware"
)

func main() {
	fmt.Println("Started")
	hardware.InitHost()
	powerbreaker := hardware.NewPowerBreaker()
	powerbreaker.SetHigh()
	time.Sleep(3 * time.Second)
	powerbreaker.SetLow()
}
