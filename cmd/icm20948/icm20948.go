package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/mpu"
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

	_ = mpu.SetAccelerometerConfig(3)

	accConfig, err := mpu.GetAccelerometerConfig()
	fmt.Println(accConfig)

	var accX, accY, accZ float64
	tstart := time.Now()
	const loops = 30000
	for i := 0; i < loops; i++ {
		accX, accY, accZ, _, _, _, err = mpu.ReadData()
	}
	elapsed := time.Since(tstart).Seconds()
	fmt.Println(float64(loops) / elapsed)

	fmt.Printf("accX: %f, accY: %f, accZ: %f\n", accX, accY, accZ)
}
