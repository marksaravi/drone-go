package icm20948

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/modules/mpu/gyroscope"
)

const (
	DPS_250  byte = 0
	DPS_500  byte = 1
	DPS_1000 byte = 2
	DPS_2000 byte = 3
)

// SetGyroConfig initialize the Gyroscope
func (dev *Device) SetGyroConfig(config *gyroscope.Config) error {
	data, err := dev.readRegister(GYRO_CONFIG_1, 1)
	config1 := data[0]
	config1 = config1 | (config.Scale << 1)
	fmt.Println(config1)
	err = dev.writeRegister(GYRO_CONFIG_1, config1)
	return err
}

// GetGyroConfig initialize the Gyroscope
func (dev *Device) GetGyroConfig() (*gyroscope.Config, error) {
	data, err := dev.readRegister(GYRO_CONFIG_1, 1)
	config1 := data[0]
	return &gyroscope.Config{
		Scale: config1,
	}, err
}
