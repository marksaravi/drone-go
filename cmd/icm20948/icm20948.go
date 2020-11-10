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

func toUint16(h, l byte) uint16 {
	fmt.Println(h, ", ", l)
	var i uint16 = uint16(h)
	i = (i << 8) & 0xFF00
	i = i | uint16(l)
	return i
}

func main() {
	// gyroConfig := gyroscope.Config{
	// 	Scale: icm20948.DPS_1000,
	// }
	var mpu mpu.MPU
	// r := make([]byte, 2)
	mpu, err := icm20948.NewRaspberryPiICM20948Driver(0, 0)
	errCheck("Initializing MPU", err)
	defer mpu.Close()
	mpu.SetDeviceConfig()
	config, err := mpu.GetDeviceConfig()
	prn("Device Config", config)
	name, id, err := mpu.WhoAmI()
	fmt.Printf("name: %s, id: 0x%X\n", name, id)
	time.Sleep(100 * time.Millisecond)
	data, err := mpu.ReadData()
	accX := toUint16(data[0], data[1])
	accY := toUint16(data[2], data[3])
	accZ := toUint16(data[4], data[5])
	gyroX := toUint16(data[6], data[7])
	gyroY := toUint16(data[8], data[9])
	gyroZ := toUint16(data[10], data[11])

	fmt.Printf("accX: %d, accY: %d, accZ: %d\n", accX, accY, accZ)
	fmt.Printf("gyroX: %d, gyroY: %d, gyroZ: %d\n", gyroX, gyroY, gyroZ)

	// mpu.ResetToDefault()
	// name, id, err := mpu.WhoAmI()
	// fmt.Printf("name: %s, id: 0x%X\n", name, id)
	// mpu.SetGyroConfig(&gyroConfig)
	// gc, err := mpu.GetGyroConfig()
	// fmt.Println(gc)

	// fmt.Println(dev.Conn.String())
	// fmt.Println("MaxTxSise: ", dev.SPI.MaxTxSize())
	// fmt.Println("CLK:  ", dev.SPI.CLK())
	// fmt.Println("MISO: ", dev.SPI.MISO())
	// fmt.Println("MOSI: ", dev.SPI.MOSI())
	// fmt.Println("CS:   ", dev.SPI.CS())

	// r, err = dev.ReadRegister(icm20948.WHO_AM_I, 1)
	// prn("Who am I", r)

	// set bank 2
	// dev.Conn.Tx([]byte{icm20948.REG_BANK_SEL, icm20948.BANK2}, nil)
	// ErrCheck("selRegisterBank", dev.selRegisterBank(2))

	// // read MOD_CTRL_USR
	// r, err = dev.ReadRegister(icm20948.MOD_CTRL_USR, 1)
	// Prn("MOD_CTRL_USR bank2", r)

	// r, err = dev.ReadRegister(icm20948.WHO_AM_I, 1)
	// Prn("Who am I", r)

	// read PWR_MGMT_1
	// powermgm1, err := dev.ReadRegister(icm20948.PWR_MGMT_1, 1)
	// prn("PWR_MGMT_1", powermgm1)
	// const powersettings byte = 0b10011111
	// err = dev.WriteRegister(icm20948.PWR_MGMT_1, powermgm1[0]&powersettings)
	// powermgm1, err = dev.ReadRegister(icm20948.PWR_MGMT_1, 1)
	// errCheck("Write", err)
	// prn("PWR_MGMT_1", powermgm1)

	// r, err = dev.WeiteRegister(icm20948.WHO_AM_I, 1)
	// Prn("Who am I", r)

	// // read PWR_MGMT_1
	// ErrCheck("selRegisterBank", dev.selRegisterBank(0))
	// err = dev.WriteRegister(icm20948.PWR_MGMT_2, 0b00000111)
	// r, err = dev.ReadRegister(icm20948.PWR_MGMT_2, 1)
	// Prn("PWR_MGMT_2 bank0", r)
	// err = dev.WriteRegister(icm20948.PWR_MGMT_2, 0b00111000)
	// r, err = dev.ReadRegister(icm20948.PWR_MGMT_2, 1)
	// Prn("PWR_MGMT_2 bank0", r)

	// const gyroconfig2 byte = 0b00001010
	// err = dev.WriteRegister(icm20948.GYRO_SMPLRT_DIV, gyroconfig2)
	// prn("SET GYRO_CONFIG_2", []byte{gyroconfig2})
	// r, err = dev.ReadRegister(icm20948.GYRO_SMPLRT_DIV, 1)
	// prn("GET GYRO_CONFIG_2", r)
}
