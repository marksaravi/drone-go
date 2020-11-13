package icm20948

import (
	"log"

	"github.com/MarkSaravi/drone-go/modules/mpu/threeaxissensore"
)

// GetGyro get accelerometer data
func (dev *Device) GetGyro() threeaxissensore.ThreeAxisSensore {
	return &(dev.gyro)
}

func (dev *Device) getGyroConfig() (GyroscopeConfig, error) {
	data, err := dev.readRegister(GYRO_CONFIG_1, 2)
	config := GyroscopeConfig{
		FullScale: int((data[0] >> 1) & 0b00000011),
	}
	dev.GetGyro().SetConfig(config)
	return config, err
}

// InitGyroscope initialise the Gyroscope
func (dev *Device) InitGyroscope() error {
	config, ok := dev.GetGyro().GetConfig().(GyroscopeConfig)
	if !ok {
		log.Fatal("Gyro config mismatch")
	}
	var config1 byte = 0
	config1 = (byte(config.FullScale) << 1) | (config1 & 0b11111001)
	err := dev.writeRegister(GYRO_CONFIG_1, config1)
	return err
}
