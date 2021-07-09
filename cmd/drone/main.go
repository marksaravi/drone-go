package main

import (
	"fmt"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	flightcontrol "github.com/MarkSaravi/drone-go/flight-control"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	appConfig := readConfigs()

	var command types.Command
	var wg sync.WaitGroup
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	commandChannel := createCommandChannel(&wg)
	imu := initiateIMU(appConfig)

	var flightStates flightcontrol.FlightStates = flightcontrol.FlightStates{
		Config: appConfig.Flight,
	}

	imu.ResetReadingTimes()

	var running = true
	for running {
		available, imuRotations, err := imu.GetRotations()
		if available {
			if err == nil {
				flightStates.Update(imuRotations)
				if udpLogger.Enabled() {
					udpLogger.Append(&flightStates)
					udpLogger.Send()
				}
			}
		}

		select {
		case command = <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				fmt.Println("COMMAND_END_PROGRAM is received, terminating services...")
				running = false
				wg.Wait()
			}
		default:
		}
	}
	imu.Close()
	dataQualityReport(imu.GetReadingQualities())
}
