// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/motors"
	"github.com/marksaravi/drone-go/devices/udplogger"
	"github.com/marksaravi/drone-go/drivers/nrf204"
	"github.com/marksaravi/drone-go/flightcontrol"
)

func main() {

	imu, imuDataPerSecond, escUpdatePerSecond := devices.NewImu()
	flightControl := flightcontrol.NewFlightControl(
		imuDataPerSecond,
		escUpdatePerSecond,
		imu,
		motors.NewESC(),
		nrf204.NewRadio(),
		udplogger.NewLogger(),
	)

	ctx, cancel := context.WithCancel(context.Background())
	var workgroup sync.WaitGroup
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
