package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

func initiateIMU(config icm20948.Config) imu.IMU {
	dev, err := icm20948.NewICM20948Driver(config)
	if err != nil {
		os.Exit(1)
	}
	if dev.InitDevice() != nil {
		os.Exit(1)
	}
	return dev
}

func createImuChannel(dataPerSecond int, config icm20948.Config, wg *sync.WaitGroup) (<-chan imu.ImuData, chan<- types.Command) {
	imuDataChannel := make(chan imu.ImuData, 64)
	imuControlChannel := make(chan types.Command, 1)
	go func() {
		wg.Add(1)
		dev := initiateIMU(config)

		name, code, deverr := dev.WhoAmI()
		if deverr != nil {
			fmt.Println("Failed to initialize IMU device with error ", deverr)
			os.Exit(1)
		}
		fmt.Printf("name: %s, id: 0x%X\n", name, code)
		var control types.Command
		readingInterval := int64(time.Second) / int64(dataPerSecond)
		firstReading := time.Now()
		nextReading := firstReading
		var total int64 = 0
		var sampleRate int = 0
		var sampleCounter int = 0
		var second = firstReading

		running := true
		for running {
			if time.Now().Before(nextReading) {
				continue
			}
			select {
			case control = <-imuControlChannel:
				if control.Command == commands.COMMAND_END_PROGRAM {
					running = false
				}
			default:
				data, err := dev.ReadData()
				if err == nil {
					total += 1
					sampleCounter += 1
					nextReading = firstReading.Add(time.Duration(total * readingInterval))
					data.TotalData = total
					data.SampleRate = sampleRate
					if time.Since(second) >= time.Second {
						second = time.Now()
						sampleRate = sampleCounter
						sampleCounter = 0
					}
					select {
					case imuDataChannel <- data:
					default:
					}
				}
			}
		}
		dev.Close()
		fmt.Println("IMU stopped.")
		wg.Done()
	}()
	return imuDataChannel, imuControlChannel
}
