package main

import (
	"fmt"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	flightcontrol "github.com/MarkSaravi/drone-go/flight-control"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	var command types.Command
	var imuData imu.ImuData
	var wg sync.WaitGroup
	var flightStates flightcontrol.FlightStates
	commandChannel := createCommandChannel(&wg)
	imuDataChannel, imuControlChannel := createImuChannel(&wg)
	var running = true

	for running {
		select {
		case command = <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				fmt.Println("COMMAND_END_PROGRAM is received, terminating services...")
				select {
				case imuControlChannel <- command:
				default:
				}
				running = false
				wg.Wait()
			}
		case imuData = <-imuDataChannel:
			flightStates.SetImuData(imuData)
			flightStates.ShowStates()
		}
	}
	fmt.Println("Program stopped.")
}
