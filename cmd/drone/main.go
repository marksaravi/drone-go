// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/drivers"
	"github.com/marksaravi/drone-go/flightcontrol"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	drivers.InitHost()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	// command, connection := radioreceiver.NewRadioReceiver(ctx, &wg)
	// logger := udplogger.NewLogger(ctx, &wg)
	imu := devices.NewImu(ctx, &wg)
	flightControl := flightcontrol.NewFlightControl(imu, nil, nil, nil)
	utils.WaitToAbortByENTER(cancel, &wg)
	flightControl.Start(ctx, &wg)
	wg.Wait()
}
