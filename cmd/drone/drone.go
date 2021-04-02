package main

import "fmt"

func main() {
	mpu := initiateMPU()

	defer mpu.Dev.Close()
	mpu.ReadData()
	fmt.Println("Orientation: ", mpu.Orientation)
}
