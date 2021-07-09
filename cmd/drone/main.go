package main

import (
	"fmt"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	flightcontrol "github.com/MarkSaravi/drone-go/flight-control"
)

func main() {
	appConfig := readConfigs()

	var wg sync.WaitGroup

	udpLogger := initUdpLogger(appConfig)
	commandChannel := createCommandChannel(&wg)
	imuChannel := createImuDataChannel(appConfig)
	pid := flightcontrol.CreatePidController()

	var running bool = true
	for running {
		select {
		case rotations := <-imuChannel:
			pid.Update(rotations)
			udpLogger.Send(rotations)

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
