package main

import (
	"fmt"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	flightcontrol "github.com/MarkSaravi/drone-go/flight-control"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	appConfig := readConfigs()

	var command types.Command
	var imuData imu.ImuData
	var wg sync.WaitGroup
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.ImuDataPerSecond)
	commandChannel := createCommandChannel(&wg)
	imuDataChannel, imuControlChannel := createImuChannel(
		appConfig.Flight.ImuDataPerSecond,
		appConfig.Devices.ICM20948,
		&wg)
	var running = true
	var flightStates flightcontrol.FlightStates = flightcontrol.FlightStates{
		Config:         appConfig.Flight,
		ImuDataChannel: imuDataChannel,
	}

	flightStates.Reset()

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
			flightStates.Set(imuData)
			json := flightStates.ImuDataToJson()
			udpLogger.Send(json)
		}
	}
	fmt.Println("Program stopped.")
}
