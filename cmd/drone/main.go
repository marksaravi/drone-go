// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/drivers"
	"github.com/marksaravi/drone-go/flightcontrol"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	drivers.InitHost()
	ctx, cancel := context.WithCancel(context.Background())
	var waitgroup sync.WaitGroup
	// command, connection := radioreceiver.NewRadioReceiver(ctx, &waitgroup)
	// logger := udplogger.NewLogger(ctx, &waitgroup)
	// imu := devices.NewImu(ctx, &waitgroup)
	flightControl := flightcontrol.NewFlightControl(nil, nil, nil, nil)
	flightControl.Start(ctx, &waitgroup)
	utils.WaitToAbortByENTER(cancel, &waitgroup)
	waitgroup.Wait()
}
