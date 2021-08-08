package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	config := utils.ReadConfigs()
	fmt.Println("Started")
	_, _, _, powerbreaker := hardware.InitHardware(config)
	powerbreaker.Connect()
	time.Sleep(2 * time.Second)
	powerbreaker.Disconnect()
}
