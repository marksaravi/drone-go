package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/mpu"
	"github.com/MarkSaravi/drone-go/utils"
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
	var prevData, currData float64 = 0, 0
	mpu.Start()
	for {
		acc, _, _ := mpu.ReadData()
		currData = utils.CalcVectorLen(acc)
		if math.Abs(currData-prevData) > 0.05 {
			fmt.Printf("accX: %f, accY: %f, accZ: %f\n", acc.X, acc.Y, acc.Z)
			prevData = currData
			// fmt.Printf("gyroX: %f, gyroY: %f, gyroZ: %f\n", gyro.X, gyro.Y, gyro.Z)
		}
	}
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

	go acc(mpu)
	time.Sleep(100 * time.Second)
}
