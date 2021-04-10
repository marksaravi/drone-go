package main

import (
	"fmt"
	"os"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	flightcontrol "github.com/MarkSaravi/drone-go/flight-control"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	appConfig, err := readConfigs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var command types.Command
	var imuData imu.ImuData
	var wg sync.WaitGroup
	var flightStates flightcontrol.FlightStates
	commandChannel := createCommandChannel(&wg)
	imuDataChannel, imuControlChannel := createImuChannel(&wg, appConfig.Devices.ICM20948)
	var running = true
	flightConfig := types.FlightConfig{
		AccLowPassFilterCoefficient:       0.05,
		RotationsLowPassFilterCoefficient: 0.05,
	}

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
			flightStates.Set(imuData, flightConfig)
			// flightStates.ShowAccStates()
			// flightStates.ShowGyroStates()
		}
	}
	fmt.Println("Program stopped.")
}
