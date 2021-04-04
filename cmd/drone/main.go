package main

import (
	"fmt"
	"time"

	commands "github.com/MarkSaravi/drone-go/constants"
	imuLib "github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	imu := initiateIMU()
	defer imu.Close()

	name, code, err := imu.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X, %v\n", name, code, err)

	var command types.Command
	var imuData imuLib.ImuData

	commandChannel := createCommandChannel()
	imuIncomingDataChannel, imuControlChannel := createImuChannel(imu)
	var counter = 0
	var start = time.Now()
	var x float64 = 0
	var y float64 = 0
	var z float64 = 0
	for command.Command != commands.COMMAND_END_PROGRAM {
		select {
		case command = <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				imuControlChannel <- command
			}
			fmt.Println("Stopping program ")
		case imuData = <-imuIncomingDataChannel:

			x = x + imuData.Gyro.Data.X*imuData.Duration
			y = y + imuData.Gyro.Data.Y*imuData.Duration
			z = imuData.Gyro.Data.Z * imuData.Duration
			counter++
			if time.Since(start) > 500*time.Millisecond {
				fmt.Println(x, y, z, counter)
				start = time.Now()
				counter = 0
			}
		}
	}
}
