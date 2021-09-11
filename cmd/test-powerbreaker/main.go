package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devicecreators"
	"github.com/MarkSaravi/drone-go/drivers"
)

func main() {
	fmt.Println("Started")
	drivers.InitHost()
	powerbreaker := devicecreators.NewPowerBreaker()
	powerbreaker.SetHigh()
	time.Sleep(3 * time.Second)
	powerbreaker.SetLow()
}
