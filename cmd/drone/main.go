// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/flightcontrol"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()

	radioNRF204 := nrf204.NewRadio()
	radioDev := radio.NewRadio(radioNRF204, 750)
	// logger := udplogger.NewLogger(&wg)
	// imudev := imu.NewImu()
	// throttles, onOff := motors.NewThrottleChannel(&wg)
	// pid := pidcontrol.NewPIDControl()
	// escRefreshInterval := time.Second / time.Duration(config.ReadFlightControlConfig().Configs.EscUpdatePerSecond)
	flightControl := flightcontrol.NewFlightControl(
		nil,
		nil,
		nil,
		nil,
		0,
		radioDev,
		nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	utils.WaitToAbortByENTER(cancel, &wg)
	radioDev.Start(ctx, &wg)
	flightControl.Start(ctx, &wg)
	log.Println("Waiting for routines to stop...")
	wg.Wait()
}
