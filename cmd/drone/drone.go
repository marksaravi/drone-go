package main

import (
	"fmt"
)

func main() {
	imu := initiateIMU()

	defer imu.Dev.Close()
	name, code, err := imu.Dev.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X, %v\n", name, code, err)
	imu.Dev.Start()
	imu.ReadData()
	fmt.Println("Orientation: ", imu.Orientation)
}
