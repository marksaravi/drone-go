package main

import (
	"fmt"
	"sync"
	"time"

	commands "github.com/MarkSaravi/drone-go/constants"
	flightcontrol "github.com/MarkSaravi/drone-go/flight-control"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	appConfig := readConfigs()

	var command types.Command
	var wg sync.WaitGroup
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.ImuDataPerSecond)
	commandChannel := createCommandChannel(&wg)
	var mpu imu.IMU = initiateIMU(appConfig.Devices.ICM20948, appConfig.Flight.LowPassFilterCoefficient)
	var running = true
	var flightStates flightcontrol.FlightStates = flightcontrol.FlightStates{
		Config: appConfig.Flight,
	}

	var readingInterval time.Duration = time.Duration(int64(time.Second) / int64(appConfig.Flight.ImuDataPerSecond))
	var badReadingInterval = readingInterval + readingInterval/10
	var max time.Duration = 0

	var prevRead = time.Now()
	var counter int64 = 0
	var badInterval int64 = 0
	var badImuCounter int64 = 0

	for running {
		now := time.Now()
		diff := now.Sub(prevRead)
		if diff < readingInterval {
			continue
		}
		counter++
		if counter%100000 == 0 {
			fmt.Println("Error Rate: ", float64(badInterval)/float64(counter)*100)
		}
		if diff >= readingInterval {
			max = diff
		}
		if diff >= badReadingInterval {
			badInterval++
		}
		prevRead = now
		imuRotations, err := mpu.GetRotations()
		if err == nil {
			flightStates.Update(imuRotations)
		} else {
			badImuCounter++
		}
		if udpLogger.Enabled() {
			udpLogger.Send(flightStates.ImuDataToJson())
		}
		select {
		case command = <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				fmt.Println("COMMAND_END_PROGRAM is received, terminating services...")
				running = false
				wg.Wait()
			}
		default:
		}
	}
	mpu.Close()
	fmt.Println("worst delay:      ", max)
	fmt.Println("total read data:  ", counter)
	fmt.Println("bad time interval:", badInterval)
	fmt.Println("bad imu data:     ", badImuCounter)
	fmt.Println("Program stopped.")
}
