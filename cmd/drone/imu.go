package main

import (
	"fmt"
	"os"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

func initiateIMU() imu.IMU {
	dev, err := icm20948.NewICM20948Driver(icm20948.Settings{
		BusNumber:  0,
		ChipSelect: 0,
		Config:     icm20948.DeviceConfig{},
		AccConfig:  icm20948.AccelerometerConfig{SensitivityLevel: 3},
		GyroConfig: icm20948.GyroscopeConfig{
			ScaleLevel:             2,
			LowPassFilterEnabled:   true,
			LowPassFilter:          7,
			LowPassFilterAveraging: 7,
			XOffset:                3.8160075,
			YOffset:                5.97310,
			ZOffset:                6.07156,
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

func createImuChannel(wg *sync.WaitGroup) (<-chan imu.ImuData, chan<- types.Command) {
	imuDataChannel := make(chan imu.ImuData)
	imuControlChannel := make(chan types.Command, 1)
	go func() {
		wg.Add(1)
		dev := initiateIMU()

		name, code, deverr := dev.WhoAmI()
		if deverr != nil {
			fmt.Println("Failed to initialize IMU device with error ", deverr)
			os.Exit(1)
		}
		fmt.Printf("name: %s, id: 0x%X\n", name, code)
		var control types.Command
		dev.ResetGyroTimer()
		running := true
		for running {
			select {
			case control = <-imuControlChannel:
				if control.Command == commands.COMMAND_END_PROGRAM {
					running = false
				}
			default:
				data, err := dev.ReadData()
				if err == nil {
					select {
					case imuDataChannel <- data:
					default:
					}
				}
			}
		}
		dev.Close()
		fmt.Println("IMU stopped.")
		wg.Done()
	}()
	return imuDataChannel, imuControlChannel
}
