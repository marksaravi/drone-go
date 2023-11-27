package main

import (
	"context"
	"fmt"
	"log"
	"sync"

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
	// log.Println("Starting Drone")
	configs := dronePackage.ReadConfigs("./configs/drone-configs.json")
	log.Println(configs)

	icm20789Configs := icm20789.Configs{
		Accelerometer: icm20789.AccelerometerConfigs{
			FullScale:              "4g",
			LowPassFilterFrequency: "44.8hz",
			NumberOfSamples:        32,
			Offsets: icm20789.Offsets{
				X: 0,
				Y: 0,
				Z: 0,
			},
		},
		Gyroscope: icm20789.GyroscopeConfigs{
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

	imuConfigs := configs.IMU
	mems := icm20789.NewICM20789(icm20789Configs)
	imudev := imu.NewIMU(mems, imu.Configs{
		DataPerSecond:   2500,
		OutputPerSecond: 5,
		AccelerometerComplimentaryFilterCoefficient: 0.02,
		RotationsComplimentaryFilterCoefficient:     0.02,
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
		for {
			var input string
			fmt.Scanln(&input)
			if input == "q" || input == "Q" {
				break
			}	
		}
		fmt.Println("Aborting Drone...")
		cancel()
	}()

	var wg sync.WaitGroup
	drone.Start(ctx, &wg)
	wg.Wait()
}
