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
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	// log.Println("Starting Drone")
	configs := dronePackage.ReadConfigs("./configs/drone-configs.json")
	log.Println(configs)

	icm20789Configs := icm20789.ReadConfigs("./configs/hardware.json")

	imuConfigs := configs.IMU
	mems := icm20789.NewICM20789(icm20789Configs)
	imudev := imu.NewIMU(mems, imu.Configs{
		DataPerSecond: 2500,
		AccelerometerComplimentaryFilterCoefficient: 0.02,
		RotationsComplimentaryFilterCoefficient:     0.02,
	})

	radioConfigs := configs.RemoteControl.Radio
	radioLink := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI,
		radioConfigs.RxTxAddress,
	)
	radioReceiver := radio.NewRadioReceiver(radioLink)

	powerBreakerGPIO := hardware.NewGPIOOutput("GPIO17")
	powerBreaker := devices.NewPowerBreaker(powerBreakerGPIO)
	b, _ := i2creg.Open("/dev/i2c-1")
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}

	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:  i2cConn,
		MaxThrottle: 15,
	})
	motorsToChannelsMappings := make(map[int]int)
	motorsToChannelsMappings[3] = 4
	motorsToChannelsMappings[1] = 6
	motorsToChannelsMappings[2] = 5
	motorsToChannelsMappings[0] = 7
	esc := esc.NewESC(pwmDev, motorsToChannelsMappings, powerBreaker, 50, false)
	ctx, cancel := context.WithCancel(context.Background())
	drone := dronePackage.NewDrone(dronePackage.DroneSettings{
		ImuDataPerSecond:  imuConfigs.DataPerSecond,
		ImuMems:           imudev,
		Escs:              esc,
		Receiver:          radioReceiver,
		CommandsPerSecond: configs.RemoteControl.CommandsPerSecond,
		PlotterActive:     configs.Plotter.Active,
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
