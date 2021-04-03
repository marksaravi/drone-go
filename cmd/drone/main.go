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

	var isRunning bool = true
	var command types.Command
	var imuData imuLib.ImuData
	var stopImuCh = make(chan bool)
	commandCh := createCommandChannel()
	imuCh := createImuChannel(imu, stopImuCh)
	for isRunning {
		select {
		case command = <-commandCh:
			if command.Command == commands.COMMAND_END_PROGRAM {
				stopImuCh <- true
				isRunning = false
			}
			fmt.Println("Stopping program ")
		case imuData = <-imuCh:
			fmt.Println(imuData)
		}
	}
}
