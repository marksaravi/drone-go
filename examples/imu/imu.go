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
		Accelerometer: icm20789.AccelerometerConfigs{
			FullScale: "4g",
			LowPassFilterFrequency: "44.8hz",
			NumberOfSamples: 32,
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
	whoAmI, err:=mems.WhoAmI()
	if err == nil {
		fmt.Printf("WHO AM I: %x\n", whoAmI)
	}
	imuConfigs := imu.Configs{
		AccelerometerComplimentaryFilterCoefficient: 0.01,
	}

	imudev := imu.NewIMU(mems, imuConfigs)

	lastRead := time.Now()

	var accrotations imu.Rotations
	running:=true
	for running {
		select {
		case <-ctx.Done():
			running=false
		default:
			if time.Since(lastRead) >= time.Second/2 {
				lastRead = time.Now()
				_, accrotations, _, _ = imudev.Read()
				fmt.Printf(" Roll: %6.2f, Pitch: %6.2f, Yaw: %6.2f\n", accrotations.Roll, accrotations.Pitch, accrotations.Yaw)
			}
		}
	}
}
