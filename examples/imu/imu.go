package main

import (
	"context"
	"fmt"
	"log"
	"time"

	plotterclient "github.com/marksaravi/drone-go/apps/plotter-client"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems/icm20789"
	"github.com/marksaravi/drone-go/timeinterval"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	plotterClient := plotterclient.NewPlotter(plotterclient.Settings{
		Active:  false,
		Address: "192.168.1.101:8000",
	})
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Scanln()
		cancel()
	}(cancel)

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

	mems := icm20789.NewICM20789(icm20789Configs)
	whoAmI, err := mems.WhoAmI()
	if err == nil {
		fmt.Printf("WHO AM I: %x\n", whoAmI)
	}
	imuConfigs := imu.Configs{
		DataPerSecond:   2500,
		OutputPerSecond: 5,
		AccelerometerComplimentaryFilterCoefficient: 0.02,
		RotationsComplimentaryFilterCoefficient:     0.02,
	}

	imudev := imu.NewIMU(mems, imuConfigs)
	running := true
	const DATA_PER_SECOND = 1000
	lastRead := time.Now()
	plotterClient.SetStartTime(lastRead)
	printInterval := timeinterval.WithDataPerSecond(4)
	readInterval := timeinterval.WithMinInterval(DATA_PER_SECOND, 25)

	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			if readInterval.IsTime() {
				rot, acc, gyro, err := imudev.Read()
				if err == nil {
					plotterClient.SendPlotterData(rot, acc, gyro)
				}
				if printInterval.IsTime() {
					printRotations(rot, acc, gyro)
				}
			}
		}
	}
}

func printRotations(rot, acc, gyro imu.Rotations) {
	// fmt.Printf("Roll: %6.2f, Pitch: %6.2f, Yaw: %6.2f,  Acc Roll: %6.2f, Pitch: %6.2f,  Gyro Roll: %6.2f, Pitch: %6.2f, Yaw: %6.2f\n", rot.Roll, rot.Pitch, rot.Yaw, acc.Roll, acc.Pitch, gyro.Roll, gyro.Pitch, gyro.Yaw)
	fmt.Printf("Roll: %6.2f, Pitch: %6.2f\n", acc.Roll, acc.Pitch)
}
