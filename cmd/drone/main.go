package main

import (
	"fmt"
	"math"
	"time"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	var command types.Command
	var imuData imu.ImuData

	commandChannel := createCommandChannel()
	imuDataChannel, imuControlChannel := createImuChannel()
	var x float64 = 0
	var y float64 = 0
	var z float64 = 0
	var currValue float64 = 1000
	lastPrint := time.Now()
	var running = true
	for running {
		select {
		case command = <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				fmt.Println("command is COMMAND_END_PROGRAM")
				select {
				case imuControlChannel <- command:
					fmt.Println("Sending COMMAND_END_PROGRAM to IMU")
				default:
					fmt.Println("No message sent for IMU")
				}
				running = false
			}

		case imuData = <-imuDataChannel:
			// x = x + imuData.Gyro.Data.X*imuData.Duration
			// y = y + imuData.Gyro.Data.Y*imuData.Duration
			// z = z + imuData.Gyro.Data.Z*imuData.Duration
			x = imuData.Acc.Data.X
			y = imuData.Acc.Data.Y
			z = imuData.Acc.Data.Z

			v := math.Sqrt(x*x + y*y + z*z)

			if math.Abs(currValue-v) > 0.025 && time.Since(lastPrint) > time.Millisecond*250 {
				// fmt.Println(x, y, z)
				// e, _ := euler.AccelerometerToEulerAngles(imuData.Acc.Data)
				// fmt.Println(fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.3f", x, y, z, e.Theta, e.Phi))
				lastPrint = time.Now()
				currValue = v
			}

		}
	}
}
