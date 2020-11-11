package icm20948

import (
	"github.com/MarkSaravi/drone-go/modules/mpu/gyroscope"
)

// SetGyroConfig initialize the Gyroscope
func (dev *Device) SetGyroConfig(config gyroscope.Config) error {
	var config1 byte = 0
	config1 = (byte(config.FullScale) << 1) | (config1 & 0b11111001)
	err := dev.writeRegister(GYRO_CONFIG_1, config1)
	return err
}

// GetGyroConfig initialize the Gyroscope
func (dev *Device) GetGyroConfig() (gyroscope.Config, error) {
	data, err := dev.readRegister(GYRO_CONFIG_1, 2)
	dev.gyroConfig.FullScale = int((data[0] >> 1) & 0b00000011)
	return dev.gyroConfig, err
}
