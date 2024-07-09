package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	dronePackage "github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/esc"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"github.com/marksaravi/drone-go/pid"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	log.Println("Starting Drone")
	configs := dronePackage.ReadConfigs("./configs/drone-configs.json")
	log.Println(configs)

	icm20789Configs := icm20789.ReadConfigs("./configs/hardware.json")

	imuConfigs := configs.IMU
	escsConfigs := configs.ESC
	mems := icm20789.NewICM20789(icm20789Configs)
	imudev := imu.NewIMU(mems, imuConfigs)

	radioLink := nrf24l01.NewNRF24L01EnhancedBurst(
		hardware.SPIConnConfigs{
			BusNumber:       configs.RemoteControl.Radio.SPI.BusNumber,
			ChipSelect:      configs.RemoteControl.Radio.SPI.ChipSelect,
			ChipEnabledGPIO: configs.RemoteControl.Radio.SPI.ChipEnabledGPIO,
		},
		configs.RemoteControl.Radio.RxTxAddress,
	)
	radioReceiver := radio.NewRadioReceiver(radioLink)

	pca9685Configs := pca9685.ReadConfigs("./configs/hardware.json")
	powerBreakerGPIO := hardware.NewGPIOOutput(pca9685Configs.BreakerGPIO)
	powerBreaker := devices.NewPowerBreaker(powerBreakerGPIO)
	b, _ := i2creg.Open(pca9685Configs.I2CPort)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}

	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:      i2cConn,
		MaxSafeThrottle: pca9685Configs.MaxSafeThrottle,
	})

	esc := esc.NewESC(pwmDev, pca9685Configs.MotorsMappings, powerBreaker, 50, false)
	ctx, cancel := context.WithCancel(context.Background())
	drone := dronePackage.NewDrone(dronePackage.DroneSettings{
		ImuDataPerSecond:  imuConfigs.DataPerSecond,
		ESCsDataPerSecond: escsConfigs.DataPerSecond,
		ImuMems:           imudev,
		Escs:              esc,
		Receiver:          radioReceiver,
		RollMidValue:      configs.Commands.RollMidValue,
		PitchMidValue:     configs.Commands.PitchMidValue,
		YawMidValue:       configs.Commands.YawMidValue,
		RotationRange:     configs.Commands.RotationRange,
		MaxThrottle:       configs.Commands.MaxThrottle,
		CommandsPerSecond: configs.RemoteControl.CommandsPerSecond,
		PlotterActive:     configs.Plotter.Active,
		PID: pid.PIDSettings{
			// MaxError:            configs.PID.MaxRotationError,
			// MaxIntegrationValue: configs.PID.MaxIntegrationValue,
			PGain: configs.PID.P,
			IGain: configs.PID.I,
			DGain: configs.PID.D,
			// MaxWeightedSum:      configs.PID.MaxWeightedSum,
			// CalibrationMode:     configs.PID.CalibrationMode,
			// CalibrationIncP:     configs.PID.CalibrationIncP,
			// CalibrationIncI:     configs.PID.CalibrationIncI,
			// CalibrationIncD:     configs.PID.CalibrationIncD,
		},
	})

	go func() {
		fmt.Scanln()
		fmt.Println("Aborting Drone...")
		cancel()
	}()

	var wg sync.WaitGroup
	drone.Start(ctx, &wg)
	wg.Wait()
}
