package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	flightcontrol "github.com/MarkSaravi/drone-go/flight-control"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/hardware/icm20948"
	"github.com/MarkSaravi/drone-go/modules/command"
	"github.com/MarkSaravi/drone-go/modules/esc"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"github.com/MarkSaravi/drone-go/types"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Flight   types.FlightConfig      `yaml:"flight_control"`
	Hardware hardware.HardwareConfig `yaml:"devices"`
	UDP      types.UdpLoggerConfig   `yaml:"udp"`
}

func main() {
	appConfig := readConfigs()
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)

	var wg sync.WaitGroup
	commandChannel := command.CreateCommandChannel(&wg)
	imu := initiateIMU(appConfig)
	pid := flightcontrol.CreatePidController()
	esc := esc.NewESCsHandler()

	var running bool = true
	imu.ResetReadingTimes()
	for running {
		if imu.CanRead() {
			rotations, err := imu.GetRotations()
			if err == nil {
				throttles := pid.Update(rotations)
				esc.SetThrottles(throttles)
				udpLogger.Send(rotations)
			}
		}
		select {
		case command := <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				fmt.Println("COMMAND_END_PROGRAM is received, terminating services...")
				wg.Wait()
				running = false
			}
		default:
		}
	}
}

func initiateIMU(config ApplicationConfig) types.IMU {
	dev, err := icm20948.NewICM20948Driver(config.Hardware.ICM20948)
	if err != nil {
		os.Exit(1)
	}
	dev.InitDevice()
	if err != nil {
		os.Exit(1)
	}
	imudevice := imu.NewIMU(dev, config.Flight.Imu)
	return &imudevice
}

func readConfigs() ApplicationConfig {
	var config ApplicationConfig

	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
		os.Exit(1)
	}
	fmt.Println(config)
	return config
}
