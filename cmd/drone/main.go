// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/apps/flightcontrol"
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/devices/udplogger"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	pidcontrol "github.com/marksaravi/drone-go/pid-control"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	configs := config.ReadConfigs().FlightControl
	radioConfigs := configs.Radio
	pidConfigs := configs.PID

	hardware.InitHost()

	radioNRF204 := nrf204.NewRadio(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)

	radioDev := radio.NewRadio(radioNRF204, radioConfigs.HeartBeatTimeoutMS)
	logger := udplogger.NewUdpLogger()
	imudev := imu.NewImu()
	pid := pidcontrol.NewPIDControl(
		pidConfigs.PGain,
		pidConfigs.IGain,
		pidConfigs.DGain,
		pidConfigs.MaxRoll,
		pidConfigs.MaxPitch,
		pidConfigs.MaxYaw,
		pidConfigs.MaxThrottle,
	)
	flightControl := flightcontrol.NewFlightControl(
		pid,
		imudev,
		radioDev,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	utils.WaitToAbortByENTER(cancel)
	radioDev.Start(&wg)
	logger.Start(&wg)
	flightControl.Start(ctx, &wg)
	wg.Wait()
}
