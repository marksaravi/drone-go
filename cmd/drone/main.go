// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/radioreceiver"
	"github.com/marksaravi/drone-go/devices/udplogger"
	"github.com/marksaravi/drone-go/drivers"
	"github.com/marksaravi/drone-go/flightcontrol"
)

func main() {
	drivers.InitHost()
	ctx, cancel := context.WithCancel(context.Background())
	var workgroup sync.WaitGroup
	command, connection := radioreceiver.NewRadioReceiver(ctx, &workgroup)
	logger := udplogger.NewLogger(ctx, &workgroup)
	imu := devices.NewImu(ctx, &workgroup)
	flightControl := flightcontrol.NewFlightControl(imu, command, connection, logger)
	workgroup.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("Press ENTER to abort")
		fmt.Scanln()
		fmt.Println("Stopping the flight control")
		cancel()
	}(&workgroup)
	flightControl.Start(ctx, &workgroup)
	workgroup.Wait()
}
