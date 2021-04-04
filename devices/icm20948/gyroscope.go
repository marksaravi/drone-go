package icm20948

import (
	"fmt"
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

	var gyroConfig1 uint8 = 0b00000000

	if config.LowPassFilterEnabled {
		gyroConfig1 = 0b00000001 | (byte(config.LowPassFilter) << 3)
	}
	gyroConfig1 = gyroConfig1 | (byte(config.ScaleLevel) << 1)

	fmt.Println("Gyro Config: ", gyroConfig1)
	err := dev.writeRegister(GYRO_CONFIG_1, gyroConfig1)
	return err
}

func (dev *Device) processGyroscopeData(data []byte) (types.XYZ, error) {
	gyroConfig, _ := dev.GetGyro().GetConfig().(GyroscopeConfig)
	gyroDegPerSec := gyroFullScale[gyroConfig.ScaleLevel]
	x := float64(utils.TowsComplementBytesToInt(data[0], data[1])) / gyroDegPerSec
	y := float64(utils.TowsComplementBytesToInt(data[2], data[3])) / gyroDegPerSec
	z := float64(utils.TowsComplementBytesToInt(data[4], data[5])) / gyroDegPerSec
	return types.XYZ{
		X: x,
		Y: y,
		Z: z,
	}, nil
}
