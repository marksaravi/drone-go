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
				X: 0,
				Y: 0,
				Z: 0,
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
	whoAmI, err:=mems.WhoAmI()
	if err == nil {
		fmt.Printf("WHO AM I: %x\n", whoAmI)
	}
	imuConfigs := imu.Configs{
		FilterCoefficient: 0.001,
	}

	imudev := imu.NewIMU(mems, imuConfigs)

	lastRead := time.Now()
	ticker := time.NewTicker(time.Second / 10)
	// var rotations, accrotations, gyrorotations imu.Rotations
	var _, accrotations, _ imu.Rotations
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// log.Printf("R: %6.2f, %6.2f, %6.2f\n",
			// 	rotations.Roll, rotations.Pitch, rotations.Yaw,
			// )
			log.Printf("A: %6.2f, %6.2f, %6.2f\n",
				accrotations.Roll, accrotations.Pitch, accrotations.Yaw,
			)
			// log.Printf("G: %6.2f, %6.2f, %6.2f\n",
			// 	gyrorotations.Roll, gyrorotations.Pitch, gyrorotations.Yaw,
			// )
		default:
			if time.Since(lastRead) >= time.Second/1000 {
				lastRead = time.Now()
				// rotations, accrotations, gyrorotations, _ = imudev.Read()
				_, accrotations, _, _ = imudev.Read()
			}
		}

	}
}
