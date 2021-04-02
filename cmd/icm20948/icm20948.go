package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/mpu"
	"github.com/MarkSaravi/drone-go/types"
)

func errCheck(step string, err error) {
	if err != nil {
		fmt.Printf("Error at %s: %s\n", step, err.Error())
		os.Exit(0)
	}
}

func readtask(dev mpu.MpuDevice, data chan types.XYZ, stop chan bool, done chan bool) {
	var gyro types.XYZ
	dev.Start()

	var finished bool = false
	for !finished {
		_, _, gyro, _, _ = dev.ReadData()
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
	dev, err := icm20948.NewICM20948Driver(icm20948.Settings{
		BusNumber:  0,
		ChipSelect: 0,
		Config:     icm20948.DeviceConfig{},
		AccConfig:  icm20948.AccelerometerConfig{SensitivityLevel: 3},
		GyroConfig: icm20948.GyroscopeConfig{ScaleLevel: 2},
	},
	)
	errCheck("Initializing ICM20948", err)
	defer dev.Close()
	dev.InitDevice()
	name, id, err := dev.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X\n", name, id)

	data := make(chan types.XYZ)
	stop := make(chan bool)
	done := make(chan bool)
	ticker := time.NewTicker(time.Second)

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Sending stop reading loop signal")
		stop <- true
	}()

	go readtask(dev, data, stop, done)

	var finished bool = false
	var counter int = 0
	var d types.XYZ
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
