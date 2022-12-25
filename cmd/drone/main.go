package main

import (
	"context"
	"log"
	"sync"

	dronepackage "github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
	"github.com/marksaravi/drone-go/types"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()

	configs := setConfigs()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	utils.WaitToAbortByESC(cancel)
	var imuMems imu.IMUMems6DOF = icm20789.NewICM20789(configs.IMU)
	var imuDevice dronepackage.InertialMeasurementUnit = imu.NewIMU(imuMems)
	drone := dronepackage.NewDrone(
		imuDevice,
	)
	drone.Start(ctx, &wg)
}

func setConfigs() types.Configs {
	return types.Configs{
		IMU: types.IMUConfigs{
			AccelerometerFullScale: "2g",
			GyroscopeFullScale:     "250dps",
		},
	}
}
