package main

import (
	"context"
	"log"
	"sync"

	dronepackage "github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	utils.WaitToAbortByESC(cancel)
	var imuMems imu.IMUMems6DOF = icm20789.NewICM20789(icm20789.ICM20789Configs{})
	var imuDevice dronepackage.InertialMeasurementUnit = imu.NewIMU(imuMems, imu.IMUConfigs{})
	drone := dronepackage.NewDrone(
		imuDevice,
	)
	drone.Start(ctx, &wg)
}
