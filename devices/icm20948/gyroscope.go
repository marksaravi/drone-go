package icm20948

import (
	"fmt"
	"log"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

// GetGyro get accelerometer data
func (dev *MemsICM20948) GetGyro() *types.Sensor {
	return &(dev.gyro)
}

func (dev *MemsICM20948) getGyroConfig() (GyroscopeConfig, error) {
	data, err := dev.readRegister(GYRO_CONFIG_1, 2)
	config := GyroscopeConfig{
		SensitivityLevel: int((data[0] >> 1) & 0b00000011),
	}
	dev.GetGyro().SetConfig(config)
	return config, err
}

// InitGyroscope initialise the Gyroscope
func (dev *MemsICM20948) InitGyroscope() error {
	config, ok := dev.GetGyro().GetConfig().(GyroscopeConfig)
	if !ok {
		log.Fatal("Gyro config mismatch")
	}

	var gyroConfig1 uint8 = 0b00000000
	var gyroConfig2 uint8 = uint8(config.Averaging)

	if config.LowPassFilterEnabled {
		gyroConfig1 = 0b00000001 | (uint8(config.LowPassFilterConfig) << 3)
	}
	gyroConfig1 = gyroConfig1 | (uint8(config.SensitivityLevel) << 1)
	err := dev.writeRegister(GYRO_CONFIG_1, gyroConfig1, gyroConfig2)
	cnfg, _ := dev.readRegister(GYRO_CONFIG_1, 2)
	fmt.Println("Gyro Config: ", cnfg)
	return err
}

func (dev *MemsICM20948) processGyroscopeData(data []uint8) (types.XYZ, error) {
	gyroConfig, _ := dev.GetGyro().GetConfig().(GyroscopeConfig)
	scale := gyroFullScale[gyroConfig.SensitivityLevel]
	offsets := gyroConfig.Offsets[gyroConfig.SensitivityLevel]
	x := (float64(utils.TowsComplementBytesToInt(data[0], data[1])) - offsets.X) / scale
	y := (float64(utils.TowsComplementBytesToInt(data[2], data[3])) - offsets.Y) / scale
	z := (float64(utils.TowsComplementBytesToInt(data[4], data[5])) - offsets.Z) / scale
	return types.XYZ{
		X: x,
		Y: y,
		Z: z,
	}, nil
}
