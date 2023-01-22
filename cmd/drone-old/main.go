// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/marksaravi/drone-go/apps/flightcontrol"
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/esc"
	imu "github.com/marksaravi/drone-go/devices/imu_old"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/devices/udplogger"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"github.com/marksaravi/drone-go/logics/pid"
	"github.com/marksaravi/drone-go/utils"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

type routine interface {
	Start(context.Context, *sync.WaitGroup)
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	flightcontrol, radioReceiver, logger := initDevices()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	utils.WaitToAbortByESC(cancel)
	radioReceiver.Start(ctx, &wg)
	logger.Start(ctx, &wg)
	flightcontrol.Start(ctx, &wg)
	wg.Wait()
	os.Exit(0)
}

type pidsettings struct {
	PGain     float64
	IGain     float64
	DGain     float64
	MaxIRatio float64
}

func createPIDSettings(
	fcConfigs pidsettings,
	maxThrottle float64,
) pid.PIDSettings {
	return pid.PIDSettings{
		PGain: fcConfigs.PGain,
		IGain: fcConfigs.IGain,
		DGain: fcConfigs.DGain,
		MaxI:  fcConfigs.MaxIRatio * maxThrottle,
	}
}

func initDevices() (flightControlRoutine, radioReceiverRoutine, udpLoggerRoutine routine) {
	configs := config.ReadConfigs()
	flightcontrolConfigs := configs.FlightControl
	fcConfigs := configs.FlightControl
	pidConfigs := fcConfigs.PID
	radioConfigs := flightcontrolConfigs.Radio
	radioNRF24L01 := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
	)
	radioReceiver := radio.NewReceiver(radioNRF24L01, flightcontrolConfigs.CommandPerSecond, radioConfigs.ConnectionTimeoutMs)
	udpLogger := udplogger.NewUdpLogger(configs.UdpLogger, flightcontrolConfigs.Imu.DataPerSecond)
	imudev := imu.NewImu(flightcontrolConfigs)
	powerBreakerPin := fcConfigs.PowerBreaker
	powerBreakerGPIO := hardware.NewGPIOOutput(powerBreakerPin)
	powerBreaker := devices.NewPowerBreaker(powerBreakerGPIO)
	b, _ := i2creg.Open(fcConfigs.ESC.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:      i2cConn,
		MaxThrottle:     fcConfigs.MaxThrottle,
		ChannelMappings: fcConfigs.ESC.PwmDeviceToESCMappings,
	})
	esc := esc.NewESC(pwmDev, powerBreaker, fcConfigs.Imu.DataPerSecond, fcConfigs.Debug)

	pidXSettings := createPIDSettings(pidsettings(pidConfigs.X), fcConfigs.MaxThrottle)
	pidYSettings := createPIDSettings(pidsettings(pidConfigs.Y), fcConfigs.MaxThrottle)
	pidZSettings := createPIDSettings(pidsettings(pidConfigs.Z), fcConfigs.MaxThrottle)

	pidcontrols := pid.NewPIDControls(
		pidXSettings,
		pidYSettings,
		pidZSettings,
		fcConfigs.Arm_0_2_ThrottleEnabled,
		fcConfigs.Arm_1_3_ThrottleEnabled,
		fcConfigs.MinPIDThrottle,
		pid.CalibrationSettings(pidConfigs.Calibration),
	)

	flightControl := flightcontrol.NewFlightControl(
		pidcontrols,
		imudev,
		esc,
		radioReceiver,
		udpLogger,
		flightcontrol.Settings{
			MaxThrottle: fcConfigs.MaxThrottle,
			MaxRoll:     fcConfigs.MaxRoll,
			MaxPitch:    fcConfigs.MaxPitch,
			MaxYaw:      fcConfigs.MaxYaw,
		},
	)
	flightControlRoutine = flightControl
	radioReceiverRoutine = radioReceiver
	udpLoggerRoutine = udpLogger
	return
}