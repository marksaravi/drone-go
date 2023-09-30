package main

import (
	"context"
	"fmt"
	"log"

	dronePackage "github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems/icm20789"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	log.Println("Starting RemoteControl")
	configs := dronePackage.ReadConfigs("./configs/drone-configs.json")
	log.Println(configs)

	imuConfigs := configs.IMU
	mems := icm20789.NewICM20789(icm20789.Configs{
		Accelerometer: icm20789.InertialDeviceConfigs{
			FullScale: imuConfigs.MEMS.Accelerometer.FullScale,
			Offsets: icm20789.Offsets{
				X: imuConfigs.MEMS.Accelerometer.Offsets.X,
				Y: imuConfigs.MEMS.Accelerometer.Offsets.Y,
				Z: imuConfigs.MEMS.Accelerometer.Offsets.Z,
			},
		},
		Gyroscope: icm20789.InertialDeviceConfigs{
			FullScale: imuConfigs.MEMS.Gyroscope.FullScale,
			Offsets: icm20789.Offsets{
				X: imuConfigs.MEMS.Gyroscope.Offsets.X,
				Y: imuConfigs.MEMS.Gyroscope.Offsets.Y,
				Z: imuConfigs.MEMS.Gyroscope.Offsets.Z,
			},
		},
		SPI: imuConfigs.MEMS.SPI,
	})
	imudev := imu.NewIMU(mems, imu.Configs{
		FilterCoefficient: imuConfigs.FilterCoefficient,
	})

	radioConfigs := configs.RemoteControl.Radio
	radioLink := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI,
		radioConfigs.RxTxAddress,
	)
	radioReceiver := radio.NewRadioReceiver(radioLink)

	ctx, cancel := context.WithCancel(context.Background())
	drone := dronePackage.NewDrone(dronePackage.DroneSettings{
		ImuDataPerSecond:  imuConfigs.DataPerSecond,
		Imu:               imudev,
		Receiver:          radioReceiver,
		CommandsPerSecond: configs.RemoteControl.CommandsPerSecond,
		PlotterActive:     configs.Plotter.Active,
	})

	go func() {
		fmt.Scanln()
		cancel()
	}()

	drone.Start(ctx)
}
