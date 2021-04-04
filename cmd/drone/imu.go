package main

import (
	"fmt"
	"os"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/imu"
	imuLib "github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

func initiateIMU() imuLib.IMU {
	dev, err := icm20948.NewICM20948Driver(icm20948.Settings{
		BusNumber:  0,
		ChipSelect: 0,
		Config:     icm20948.DeviceConfig{},
		AccConfig:  icm20948.AccelerometerConfig{SensitivityLevel: 3},
		GyroConfig: icm20948.GyroscopeConfig{
			ScaleLevel:             2,
			LowPassFilterEnabled:   true,
			LowPassFilter:          7,
			LowPassFilterAveraging: 2,
			XOffset:                0,
			YOffset:                0,
			ZOffset:                6.075,
		},
	})
	if err != nil {
		os.Exit(1)
	}
	if dev.InitDevice() != nil {
		os.Exit(1)
	}
	return dev
}

func createImuChannel(imu imuLib.IMU) (chan imu.ImuData, chan types.Command) {
	imuChannel := make(chan imuLib.ImuData)
	imuControlChannel := make(chan types.Command)
	var control types.Command
	go func() {

		imu.ResetGyroTimer()
		for control.Command != commands.COMMAND_END_PROGRAM {
			select {
			case control = <-imuControlChannel:
				if control.Command == commands.COMMAND_END_PROGRAM {
					fmt.Println("Stopping IMU")
				}
			default:
				data, err := imu.ReadData()
				if err == nil {
					imuChannel <- data
				}
			}
		}
	}()
	return imuChannel, imuControlChannel
}
