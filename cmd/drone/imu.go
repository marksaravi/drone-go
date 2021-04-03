package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/imu"
	imuLib "github.com/MarkSaravi/drone-go/modules/imu"
)

func initiateIMU() imuLib.IMU {
	dev, err := icm20948.NewICM20948Driver(icm20948.Settings{
		BusNumber:  0,
		ChipSelect: 0,
		Config:     icm20948.DeviceConfig{},
		AccConfig:  icm20948.AccelerometerConfig{SensitivityLevel: 3},
		GyroConfig: icm20948.GyroscopeConfig{ScaleLevel: 2},
	},
	)
	if err != nil {
		os.Exit(1)
	}
	if dev.InitDevice() != nil {
		os.Exit(1)
	}
	return dev
}

func createImuChannel(imu imuLib.IMU, stopImuCh <-chan bool) chan imu.ImuData {
	imuCh := make(chan imuLib.ImuData)
	var stop bool = false
	go func() {
		for !stop {
			select {
			case stop = <-stopImuCh:
				fmt.Println("Stopping IMU")
			default:
				data, err := imu.ReadData()
				if err == nil {
					imuCh <- data
				}
				time.Sleep(time.Second)
			}
		}
	}()
	return imuCh
}
