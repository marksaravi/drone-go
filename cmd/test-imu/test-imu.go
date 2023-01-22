package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
	"github.com/marksaravi/drone-go/types"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Scanln()
		cancel()
	}(cancel)

	icm20789Configs := icm20789.ICM20789Configs{
		AccelerometerFullScale: "2g",
		GyroscopeFullScale:     "250dps",
	}

	imuConfigs := imu.IMUConfigs{
		FilterCoefficient: 0.01,
	}

	dev := icm20789.NewICM20789(icm20789Configs)

	imu := imu.NewIMU(dev, imuConfigs)

	lastRead := time.Now()
	ticker := time.NewTicker(time.Second / 2)
	var rotations types.Rotations
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Printf("%6.2f, %6.2f, %6.2f\n", rotations.Roll, rotations.Pitch, rotations.Yaw)
		default:
			if time.Since(lastRead) >= time.Millisecond {
				lastRead = time.Now()
				rotations, _ = imu.Read()
			}
		}

	}
}
