package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/mpu"
	"github.com/MarkSaravi/drone-go/modules/mpu/threeaxissensore"
)

func errCheck(step string, err error) {
	if err != nil {
		fmt.Printf("Error at %s: %s\n", step, err.Error())
		os.Exit(0)
	}
}

func readtask(mpu mpu.MPU, data chan threeaxissensore.Data, stop chan bool, done chan bool) {
	var gyro threeaxissensore.Data
	mpu.Start()

	var finished bool = false
	for !finished {
		_, gyro, _ = mpu.ReadData()
		data <- gyro
		select {
		case finished = <-stop:
		default:
		}
	}
	fmt.Println("Reading loop stopped")
	fmt.Println("Sending stop programm signal")
	done <- true
}

func main() {
	var mpu mpu.MPU
	mpu, err := icm20948.NewRaspberryPiICM20948Driver(
		0,
		0,
		icm20948.DeviceConfig{},
		icm20948.AccelerometerConfig{Sensitivity: 3},
		icm20948.GyroscopeConfig{FullScale: 2},
	)
	errCheck("Initializing MPU", err)
	defer mpu.Close()
	mpu.InitDevice()
	name, id, err := mpu.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X\n", name, id)
	config, accConfig, gyroConfig, err := mpu.GetDeviceConfig()
	fmt.Println(config)
	fmt.Println(accConfig)
	fmt.Println(gyroConfig)

	data := make(chan threeaxissensore.Data)
	stop := make(chan bool)
	done := make(chan bool)
	ticker := time.NewTicker(time.Second)

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Sending stop reading loop signal")
		stop <- true
	}()

	go readtask(mpu, data, stop, done)

	var finished bool = false
	var counter int = 0
	var d threeaxissensore.Data
	for !finished {
		select {
		case finished = <-done:
			fmt.Println("Stopping program")
		case d = <-data:
			counter++
		case <-ticker.C:
			fmt.Println(d, counter)
			counter = 0
		}
	}
}
