package main

import (
	"context"
	"log"
	"sync"

	dronepackage "github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware/icm20789"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	utils.WaitToAbortByESC(cancel)
	drone := dronepackage.NewDrone(
		imu.NewIMU(icm20789.NewICM20789()),
	)
	drone.Start(ctx, &wg)
}
