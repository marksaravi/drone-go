package icm20948

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/modules/mpu/accelerometer"
)

// GetAccelerometerConfig reads Accelerometer's settings
func (dev *Device) GetAccelerometerConfig() ([]byte, error) {
	config, err := dev.readRegister(ACCEL_CONFIG, 2)
	dev.accelerometerConfig.Sensitivity = int((config[0] >> 1) & 0b00000011)
	return config, err
}

// SetAccelerometerConfig reads Accelerometer's settings
func (dev *Device) SetAccelerometerConfig(config accelerometer.Config) error {
	fmt.Println("accelerometerSensitivity", accelerometerSensitivity)
	dev.accelerometerConfig.Sensitivity = config.Sensitivity
	data, err := dev.readRegister(ACCEL_CONFIG, 2)
	var accsen byte = byte(dev.accelerometerConfig.Sensitivity) << 1
	data[0] = data[0] & 0b11111001
	data[0] = data[0] | accsen
	err = dev.writeRegister(ACCEL_CONFIG, data[0], data[1])
	time.Sleep(time.Millisecond * 100)
	return err
}
