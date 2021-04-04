package main

import (
	"fmt"

	commands "github.com/MarkSaravi/drone-go/constants"
	imuLib "github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	imu := initiateIMU()
	defer imu.Close()

	name, code, err := imu.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X, %v\n", name, code, err)
	imu.Start()

	var command types.Command
	var imuData imuLib.ImuData

	commandChannel := createCommandChannel()
	imuIncomingDataChannel, imuControlChannel := createImuChannel(imu)
	for command.Command != commands.COMMAND_END_PROGRAM {
		select {
		case command = <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				imuControlChannel <- command
			}
			fmt.Println("Stopping program ")
		case imuData = <-imuIncomingDataChannel:
			fmt.Println(imuData)
		}
	}
}
