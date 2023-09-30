package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems/icm20789"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Scanln()
		cancel()
	}(cancel)
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
		DataPerSecond:     1000,
		FilterCoefficient: 0.001,
	}

	imudev := imu.NewIMU(mems, imuConfigs)

	lastRead := time.Now()
	ticker := time.NewTicker(time.Second / 10)
	var rotations imu.Rotations
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Printf("%6.2f, %6.2f, %6.2f\n", rotations.Roll, rotations.Pitch, rotations.Yaw)
		default:
			if time.Since(lastRead) >= time.Second/time.Duration(imuConfigs.DataPerSecond) {
				lastRead = time.Now()
				rotations, _ = imudev.Read()
			}
		}

	}
}
