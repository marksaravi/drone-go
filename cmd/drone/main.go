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
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.ImuDataPerSecond)
	commandChannel := createCommandChannel(&wg)
	imu := initiateIMU(appConfig)
	var running = true
	var flightStates flightcontrol.FlightStates = flightcontrol.FlightStates{
		Config: appConfig.Flight,
	}

	imu.ResetReadingTimes()
	for running {
		available, imuRotations, err := imu.GetRotations()
		if available {
			if err == nil {
				flightStates.Update(imuRotations)
				if udpLogger.Enabled() {
					udpLogger.Append(flightStates.ImuDataToJson())
				}
			}
		}
		udpLogger.Send()
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
	readingData := imu.GetReadingQualities()
	fmt.Println("total data:             ", readingData.Total)
	fmt.Println("number of bad imu data: ", readingData.BadData)
	fmt.Println("number of bad timing:   ", readingData.BadInterval)
	fmt.Println("bad timing rate:        ", float64(readingData.BadInterval)/float64(readingData.Total)*100)
	fmt.Println("max bad timing:         ", readingData.MaxBadInterval)
	fmt.Println("Program stopped.")
}
