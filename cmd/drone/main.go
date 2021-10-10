// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/devices/motors"
	"github.com/marksaravi/drone-go/devices/radioreceiver"
	"github.com/marksaravi/drone-go/devices/udplogger"
	"github.com/marksaravi/drone-go/flightcontrol"
	"github.com/marksaravi/drone-go/hardware"
	pidcontrol "github.com/marksaravi/drone-go/pid-control"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	command, connection := radioreceiver.NewRadioReceiver(ctx, &wg)
	logger := udplogger.NewLogger(&wg)
	imudev := imu.NewImu()
	throttles, onOff := motors.NewThrottleChannel(&wg)
	pid := pidcontrol.NewPIDControl()
	flightControl := flightcontrol.NewFlightControl(pid, imudev, throttles, onOff, command, connection, logger)
	utils.WaitToAbortByENTER(cancel, &wg)
	flightControl.Start(&wg)
	wg.Wait()
}
