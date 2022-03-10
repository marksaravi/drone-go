// Drone is the main application to run the FlightControl.
package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/apps/flightcontrol"
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/esc"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/devices/udplogger"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"github.com/marksaravi/drone-go/logics/pid"
	"github.com/marksaravi/drone-go/utils"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	configs := config.ReadConfigs()
	flightcontrolConfigs := configs.FlightControl
	fcConfigs := configs.FlightControl
	pidConfigs := fcConfigs.PID

	hardware.InitHost()

	radioReceiver := createRadioReceiver(flightcontrolConfigs)
	logger := udplogger.NewUdpLogger()
	imudev := imu.NewImu()
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

	pidRollSettings := createPIDSettings(pidsettings(pidConfigs.Roll), fcConfigs.MaxThrottle)
	pidPitchSettings := createPIDSettings(pidsettings(pidConfigs.Pitch), fcConfigs.MaxThrottle)
	pidYawSettings := createPIDSettings(pidsettings(pidConfigs.Yaw), fcConfigs.MaxThrottle)

	pidcontrols := pid.NewPIDControls(
		pidRollSettings,
		pidPitchSettings,
		pidYawSettings,
		fcConfigs.Arm_0_2_ThrottleEnabled,
		fcConfigs.Arm_1_3_ThrottleEnabled,
		fcConfigs.MinPIDThrottle,
		pid.CalibrationSettings(pidConfigs.Calibration),
	)
	fmt.Println(pidcontrols)
	flightControl := flightcontrol.NewFlightControl(
		pidcontrols,
		imudev,
		esc,
		radioReceiver,
		logger,
		flightcontrol.Settings{
			MaxThrottle: fcConfigs.MaxThrottle,
			MaxRoll:     fcConfigs.MaxRoll,
			MaxPitch:    fcConfigs.MaxPitch,
			MaxYaw:      fcConfigs.MaxYaw,
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	utils.WaitToAbortByENTER(cancel, &wg)
	radioReceiver.StartReceiver(ctx, &wg)
	logger.Start(&wg)
	flightControl.Start(ctx, &wg)
	wg.Wait()
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
