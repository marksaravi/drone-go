package main

import (
	"fmt"
	"math"

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
	var x float64 = 0
	var y float64 = 0
	var z float64 = 0
	var currZ float64 = 0
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
			z = z + imuData.Gyro.Data.Z*imuData.Duration
			if math.Abs(currZ-z) > 1 {
				fmt.Println(z)
				currZ = z
			}
		}
	}
}
