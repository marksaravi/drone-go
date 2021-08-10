package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware"
)

func main() {
	config := config.ReadConfigs()
	fmt.Println("Started")
	_, _, _, powerbreaker := hardware.InitDroneHardware(config)
	powerbreaker.Connect()
	time.Sleep(2 * time.Second)
	powerbreaker.Disconnect()
}
