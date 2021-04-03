package icm20948

import (
	"log"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

// GetGyro get accelerometer data
func (dev *Device) GetGyro() *types.Sensor {
	return &(dev.gyro)
}

func (dev *Device) getGyroConfig() (GyroscopeConfig, error) {
	data, err := dev.readRegister(GYRO_CONFIG_1, 2)
	config := GyroscopeConfig{
		ScaleLevel: int((data[0] >> 1) & 0b00000011),
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
	config1 = (byte(config.ScaleLevel) << 1) | (config1 & 0b11111001)
	err := dev.writeRegister(GYRO_CONFIG_1, config1)
	return err
}

func (dev *Device) processGyroscopeData(data []byte) (types.XYZ, error) {
	gyroConfig, _ := dev.GetGyro().GetConfig().(GyroscopeConfig)
	gyroDegPerSec := gyroFullScale[gyroConfig.ScaleLevel]
	x := float64(utils.TowsComplementBytesToInt(data[0], data[1])) / gyroDegPerSec
	y := float64(utils.TowsComplementBytesToInt(data[2], data[3])) / gyroDegPerSec
	z := float64(utils.TowsComplementBytesToInt(data[4], data[5])) / gyroDegPerSec
	dx := float64(dev.duration) * x / 1e9
	dy := float64(dev.duration) * y / 1e9
	dz := float64(dev.duration) * z / 1e9
	return types.XYZ{
		X: dx,
		Y: dy,
		Z: dz,
	}, nil
}
