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
	pidcontrol "github.com/marksaravi/drone-go/pid-control"
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
	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Address, i2cConn, configs.ESC.SafetyMaxThrottle)
	esc := esc.NewESC(pwmDev, powerBreaker, configs.Imu.DataPerSecond, configs.ESC.PwmDeviceToESCMappings, configs.Debug)

	pid := pidcontrol.NewPIDControl(
		configs.Imu.DataPerSecond,
		pidConfigs.PGain,
		pidConfigs.IGain,
		pidConfigs.DGain,
		pidConfigs.MaxRoll,
		pidConfigs.MaxPitch,
		pidConfigs.MaxYaw,
		pidConfigs.MaxThrottle,
		mcp3008.DIGITAL_MAX_VALUE,
	)
	flightControl := flightcontrol.NewFlightControl(
		pid,
		imudev,
		esc,
		radioDev,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	utils.WaitToAbortByENTER(cancel)
	radioDev.Start(ctx, &wg)
	logger.Start(&wg)
	esc.Start(&wg)
	flightControl.Start(ctx, &wg)
	wg.Wait()
}
