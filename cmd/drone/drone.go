package main

import (
	"fmt"
)

func main() {
	mpu := initiateMPU()

	defer mpu.Dev.Close()
	name, code, err := mpu.Dev.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X, %v\n", name, code, err)
	mpu.ReadData()
	fmt.Println("Orientation: ", mpu.Orientation)
}
