package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/mpu"
	"github.com/MarkSaravi/drone-go/modules/mpu/accelerometer"
)

func errCheck(step string, err error) {
	if err != nil {
		fmt.Printf("Error at %s: %s\n", step, err.Error())
		os.Exit(0)
	}
}

func prn(msg string, bytes []byte) {
	fmt.Printf("%s: ", msg)
	for _, b := range bytes {
		fmt.Printf("0x%X, ", b)
	}
	fmt.Printf("\n")
}

func acc(mpu mpu.MPU) {
	var accX, accY, accZ, lastvector, newvector, diff float64 = 0, 0, 0, 0, 0, 0
	for {
		accX, accY, accZ, _, _, _, _ = mpu.ReadData()
		newvector = math.Sqrt(accX*accX + accY*accY + accZ*accZ)
		diff = math.Abs(newvector - lastvector)
		if diff > 0.05 {
			fmt.Printf("accX: %f, accY: %f, accZ: %f\n", accX, accY, accZ)
			lastvector = newvector
		}
	}
}

func main() {
	var mpu mpu.MPU
	mpu, err := icm20948.NewRaspberryPiICM20948Driver(0, 0)
	errCheck("Initializing MPU", err)
	defer mpu.Close()
	mpu.SetDeviceConfig()
	config, err := mpu.GetDeviceConfig()
	prn("Device Config", config)
	name, id, err := mpu.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X\n", name, id)

	_ = mpu.SetAccelerometerConfig(accelerometer.Config{Sensitivity: 3})

	accConfig, err := mpu.GetAccelerometerConfig()
	fmt.Println(accConfig)

	go acc(mpu)
	time.Sleep(20 * time.Second)
}
