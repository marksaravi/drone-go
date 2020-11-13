package icm20948

import (
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/modules/mpu/threeaxissensore"
)

// GetAcc get accelerometer data
func (dev *Device) GetAcc() threeaxissensore.ThreeAxisSensore {
	return &(dev.acc)
}

func (dev *Device) getAccConfig() (AccelerometerConfig, error) {
	data, err := dev.readRegister(ACCEL_CONFIG, 2)
	config := AccelerometerConfig{
		Sensitivity: int((data[0] >> 1) & 0b00000011),
	}
	dev.GetAcc().SetConfig(config)
	return config, err
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
