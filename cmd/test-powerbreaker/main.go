package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	fmt.Println("Started")
	drivers.InitHost()
	powerbreaker := utils.NewPowerBreaker()
	powerbreaker.SetHigh()
	time.Sleep(3 * time.Second)
	powerbreaker.SetLow()
}
