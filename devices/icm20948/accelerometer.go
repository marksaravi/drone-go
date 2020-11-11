package icm20948

import (
	"fmt"
	"time"
)

// GetAccelerometerConfig reads Accelerometer's settings
func (dev *Device) GetAccelerometerConfig() ([]byte, error) {
	config, err := dev.readRegister(ACCEL_CONFIG, 2)
	dev.accelerometerSensitivity = int((config[0] >> 1) & 0b00000011)
	return config, err
}

// SetAccelerometerConfig reads Accelerometer's settings
func (dev *Device) SetAccelerometerConfig(accelerometerSensitivity int) error {
	fmt.Println("accelerometerSensitivity", accelerometerSensitivity)
	dev.accelerometerSensitivity = accelerometerSensitivity
	config, err := dev.readRegister(ACCEL_CONFIG, 2)
	var accsen byte = byte(accelerometerSensitivity) << 1
	config[0] = config[0] & 0b11111001
	config[0] = config[0] | accsen
	err = dev.writeRegister(ACCEL_CONFIG, config[0], config[1])
	time.Sleep(time.Millisecond * 100)
	return err
}
