// Drone is the main application to run the FlightControl.
package main

import (
	"github.com/MarkSaravi/drone-go/devicecreators"
	"github.com/MarkSaravi/drone-go/flightcontrol"
)

func main() {

	imu, imuDataPerSecond, escUpdatePerSecond := devicecreators.NewImu()
	flightControl := flightcontrol.NewFlightControl(
		imuDataPerSecond,
		escUpdatePerSecond,
		imu,
		devicecreators.NewESC(),
		devicecreators.NewRadio(),
		devicecreators.NewLogger(),
	)

	flightControl.Start()
}
