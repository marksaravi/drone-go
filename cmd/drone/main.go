package main

import (
	"fmt"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	appConfig := readConfigs()

	var command types.Command
	var wg sync.WaitGroup
	udpLogger := initUdpLogger(appConfig)
	commandChannel := createCommandChannel(&wg)
	imu := initiateIMU(appConfig)

	imu.ResetReadingTimes()

	var running = true
	for running {
		imuRotations, err := imu.GetRotations()
		if err == nil {
			udpLogger.Send(imuRotations)
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
