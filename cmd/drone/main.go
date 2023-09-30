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

	radioConfigs := configs.Radio
	radioLink := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI,
		radioConfigs.RxTxAddress,
	)
	radioReceiver := radio.NewRadioReceiver(radioLink)
	icm20789Configs := icm20789.Configs{
		Accelerometer: icm20789.InertialDeviceConfigs{
			FullScale: "4g",
			Offsets: icm20789.Offsets{
				X: 32696,
				Y: 836,
				Z: 31979,
			},
		},
		Gyroscope: icm20789.InertialDeviceConfigs{
			FullScale: "500dps",
			Offsets: icm20789.Offsets{
				X: 0,
				Y: 0,
				Z: 0,
			},
		},
		SPI: hardware.SPIConnConfigs{
			BusNumber:  0,
			ChipSelect: 0,
		},
	}

	mems := icm20789.NewICM20789(icm20789Configs)
	imuConfigs := imu.Configs{
		FilterCoefficient: 0.001,
	}

	imudev := imu.NewIMU(mems, imuConfigs)
	ctx, cancel := context.WithCancel(context.Background())
	drone := dronePackage.NewDrone(dronePackage.DroneSettings{
		ImuDataPerSecond:  1000,
		Imu:               imudev,
		Receiver:          radioReceiver,
		CommandsPerSecond: configs.CommandsPerSecond,
	})

	go func() {
		fmt.Scanln()
		cancel()
	}()

	drone.Start(ctx)
}
