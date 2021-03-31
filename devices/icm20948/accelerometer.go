package icm20948

import (
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/types/sensore"
	"github.com/MarkSaravi/drone-go/utils"
)

// GetAcc get accelerometer data
func (dev *Device) GetAcc() sensore.ThreeAxisSensore {
	return &(dev.acc)
}

// InitAccelerometer initialise the Accelerometer
func (dev *Device) InitAccelerometer() error {
	config, ok := dev.GetAcc().GetConfig().(AccelerometerConfig)
	if !ok {
		log.Fatal("Accelerometer config mismatch")
	}

	fmt.Println("Accelerometer Sensitivity", config.Sensitivity)
	data, err := dev.readRegister(ACCEL_CONFIG, 2)
	var accsen byte = byte(config.Sensitivity) << 1
	data[0] = data[0] & 0b11111001
	data[0] = data[0] | accsen
	err = dev.writeRegister(ACCEL_CONFIG, data[0], data[1])
	time.Sleep(time.Millisecond * 100)
	return err
}

func (dev *Device) getAccConfig() (AccelerometerConfig, error) {
	data, err := dev.readRegister(ACCEL_CONFIG, 2)
	config := AccelerometerConfig{
		Sensitivity: int((data[0] >> 1) & 0b00000011),
	}
	dev.GetAcc().SetConfig(config)
	return config, err
}

func (dev *Device) processAccelerometerData(data []byte) {
	accConfig, _ := dev.GetAcc().GetConfig().(AccelerometerConfig)
	accSens := accelerometerSensitivity[accConfig.Sensitivity]
	x := float64(utils.TowsComplementBytesToInt(data[0], data[1])) / accSens
	y := float64(utils.TowsComplementBytesToInt(data[2], data[3])) / accSens
	z := float64(utils.TowsComplementBytesToInt(data[4], data[5])) / accSens
	dev.GetAcc().SetData(x, y, z)
}
