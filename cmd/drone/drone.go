package main

import (
	"fmt"

	imuLib "github.com/MarkSaravi/drone-go/modules/imu"
)

func main() {
	var imu imuLib.IMU
	imu = initiateIMU()

	defer imu.Close()
	name, code, err := imu.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X, %v\n", name, code, err)
	imu.Start()
	imuData, _ := imu.ReadData()
	fmt.Println("Orientation: ", imuData.Gyro)
}
