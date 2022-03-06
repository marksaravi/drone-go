// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/apps/flightcontrol"
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/esc"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/devices/udplogger"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"github.com/marksaravi/drone-go/logics/pid"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	configs := config.ReadConfigs().FlightControl
	radioConfigs := configs.Radio
	pidConfigs := configs.PID

	hardware.InitHost()

	radioNRF204 := nrf204.NewNRF204EnhancedBurst(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
	)

	radioDev := radio.NewReceiver(radioNRF204, configs.CommandPerSecond, radioConfigs.ConnectionTimeoutMs)
	logger := udplogger.NewUdpLogger()
	imudev := imu.NewImu()
	powerBreakerPin := configs.PowerBreaker
	powerBreakerGPIO := hardware.NewGPIOOutput(powerBreakerPin)
	powerBreaker := devices.NewPowerBreaker(powerBreakerGPIO)
	b, _ := i2creg.Open(configs.ESC.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:      i2cConn,
		MaxThrottle:     configs.MaxThrottle,
		ChannelMappings: configs.ESC.PwmDeviceToESCMappings,
	})
	esc := esc.NewESC(pwmDev, powerBreaker, configs.Imu.DataPerSecond, configs.Debug)

	pidcontrols := pid.NewPIDControls(
		pid.PIDControlSettings{
			Roll:        models.PIDSettings(pidConfigs.Roll),
			Pitch:       models.PIDSettings(pidConfigs.Pitch),
			Yaw:         models.PIDSettings(pidConfigs.Yaw),
			Calibration: pid.CalibrationSettings(pidConfigs.Calibration),
		},
	)
	flightControl := flightcontrol.NewFlightControl(
		pidcontrols,
		imudev,
		esc,
		radioDev,
		logger,
		flightcontrol.Settings{
			MaxThrottle: configs.MaxThrottle,
			MaxRoll:     configs.MaxRoll,
			MaxPitch:    configs.MaxPitch,
			MaxYaw:      configs.MaxYaw,
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	utils.WaitToAbortByENTER(cancel, &wg)
	radioDev.StartReceiver(ctx, &wg)
	logger.Start(&wg)
	flightControl.Start(ctx, &wg)
	wg.Wait()
}
