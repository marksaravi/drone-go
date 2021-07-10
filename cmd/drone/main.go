package main

import (
	"fmt"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	flightcontrol "github.com/MarkSaravi/drone-go/flight-control"
	"github.com/MarkSaravi/drone-go/modules/command"
)

func main() {
	appConfig := readConfigs()

	var wg sync.WaitGroup

	udpLogger := initUdpLogger(appConfig)
	commandChannel := command.CreateCommandChannel(&wg)
	imu := initiateIMU(appConfig)
	pid := flightcontrol.CreatePidController()
	esc := createESCsHandler()

	var running bool = true
	imu.ResetReadingTimes()
	for running {
		if imu.CanRead() {
			rotations, err := imu.GetRotations()
			if err == nil {
				throttles := pid.Update(rotations)
				esc.SetThrottles(throttles)
				udpLogger.Send(rotations)
			}
		}
		select {
		case command := <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				fmt.Println("COMMAND_END_PROGRAM is received, terminating services...")
				wg.Wait()
				running = false
			}
		default:
		}
	}
}
