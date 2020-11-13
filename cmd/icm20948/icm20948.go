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
	ticker := time.NewTicker(time.Second)
	mpu.Start()
	var finished bool = false
	for !finished {
		_, gyro, _ = mpu.ReadData()
		select {
		case finished = <-stop:
			ticker.Stop()
			fmt.Println("Stopping ticker")
		case <-ticker.C:
			data <- gyro
		}
	}
	fmt.Println("Stopping data reader loop")
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

	stop := make(chan bool)
	done := make(chan bool)
	data := make(chan threeaxissensore.Data)

	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("Sending Stop Singnal")
		stop <- true
	}()

	go readtask(mpu, data, stop, done)

	var finished bool = false
	for !finished {
		select {
		case finished = <-done:
			fmt.Println("Stopping program")
		case d := <-data:
			fmt.Println(d)
		}
	}
}
