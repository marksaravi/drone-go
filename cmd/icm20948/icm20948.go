package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/mpu"
)

func main() {
	mpu, err := mpu.NewMPU("/dev/spidev0.0")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer mpu.Close()

	b, _ := mpu.GetRegister(icm20948.WHO_AM_I, icm20948.BANK0)
	fmt.Printf("WHO_AM_I: 0x%X\n", b)

	b, _ = mpu.GetRegister(icm20948.MOD_CTRL_USR, icm20948.BANK2)
	fmt.Printf("MOD_CTRL_USR: 0x%X\n", b)

	b, _ = mpu.GetRegister(icm20948.GYRO_CONFIG_1, icm20948.BANK2)
	fmt.Printf("GYRO_CONFIG_1: 0x%X\n", b)
	mpu.SetRegister(icm20948.GYRO_CONFIG_1, icm20948.BANK2, 0b01111111)
	fmt.Printf("GYRO_CONFIG_1: 0x%X\n", b)
	b, _ = mpu.GetRegister(icm20948.GYRO_CONFIG_1, icm20948.BANK2)
	fmt.Printf("GYRO_CONFIG_1: 0x%X\n", b)

	mpu.SetRegister(icm20948.GYRO_SMPLRT_DIV, icm20948.BANK2, 0x17)
	b, _ = mpu.GetRegister(icm20948.GYRO_SMPLRT_DIV, icm20948.BANK2)
	fmt.Printf("GYRO_SMPLRT_DIV: 0x%X\n", b)
}
