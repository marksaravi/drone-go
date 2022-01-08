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
	"github.com/marksaravi/drone-go/hardware/mcp3008"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"github.com/marksaravi/drone-go/pid"
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

	radioNRF204 := nrf204.NewNRF204(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)

	radioDev := radio.NewRadio(radioNRF204, radioConfigs.HeartBeatTimeoutMS)
	logger := udplogger.NewUdpLogger()
	imudev := imu.NewImu()
	powerBreakerPin := configs.PowerBreaker
	powerBreakerGPIO := hardware.NewGPIOOutput(powerBreakerPin)
	powerBreaker := devices.NewPowerBreaker(powerBreakerGPIO)
	b, _ := i2creg.Open(configs.ESC.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:           i2cConn,
		SafeStartThrottle:    configs.SafeStartThrottle,
		MaxThrottle:          configs.MaxThrottle,
		ControlVariableRange: configs.ControlVariableRange,
		ChannelMappings:      configs.ESC.PwmDeviceToESCMappings,
	})
	esc := esc.NewESC(pwmDev, powerBreaker, configs.Imu.DataPerSecond, configs.Debug)

	pidcontrols := pid.NewPIDControls(
		pid.PIDSettings{
			RollPitchPGain:          pidConfigs.RollPitchPGain,
			RollPitchIGain:          pidConfigs.RollPitchIGain,
			RollPitchDGain:          pidConfigs.RollPitchDGain,
			YawPGain:                pidConfigs.YawPGain,
			YawIGain:                pidConfigs.YawIGain,
			YawDGain:                pidConfigs.YawDGain,
			LimitRoll:               pidConfigs.MaxRoll,
			LimitPitch:              pidConfigs.MaxPitch,
			LimitYaw:                pidConfigs.MaxYaw,
			LimitI:                  pidConfigs.MaxI,
			ThrottleLimit:           float64(configs.MaxThrottle),
			SafeStartThrottle:       float64(configs.SafeStartThrottle),
			MaxJoystickDigitalValue: mcp3008.DIGITAL_MAX_VALUE,
			BeamToAxisRatio:         pidConfigs.BeamToAxisRatio,
			CalibrationGain:         pidConfigs.CalibrationGain,
			CalibrationStep:         pidConfigs.CalibrationStep,
		},
	)
	flightControl := flightcontrol.NewFlightControl(
		pidcontrols,
		imudev,
		esc,
		radioDev,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	utils.WaitToAbortByENTER(cancel)
	radioDev.Start(ctx, &wg)
	// logger.Start(&wg)
	esc.Start(&wg)
	flightControl.Start(ctx, &wg)
	wg.Wait()
}
