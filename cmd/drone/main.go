package main

import (
	"context"
	"log"
	"sync"

	dronepackage "github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems/icm20789"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	utils.WaitToAbortByESC(cancel)
	var mems imu.IMUMems6DOF = icm20789.NewICM20789(icm20789.ReadConfigs())
	imudev := imu.NewIMU(mems, imu.ReadConfigs())
	drone := dronepackage.NewDrone(
		imudev,
	)
	drone.Fly(ctx, &wg)
}
