package icm20948

import (
	"fmt"
	"log"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

// GetGyro get accelerometer data
func (dev *memsICM20948) GetGyro() *types.Sensor {
	return &(dev.gyro)
}

// InitGyroscope initialise the Gyroscope
func (dev *memsICM20948) InitGyroscope() error {
	config, ok := dev.GetGyro().GetConfig().(types.GyroscopeConfig)
	if !ok {
		log.Fatal("Gyro config mismatch")
	}

	var gyroConfig1 uint8 = 0b00000000
	var gyroConfig2 uint8 = uint8(config.Averaging)

	if config.LowPassFilterEnabled {
		gyroConfig1 = 0b00000001 | (uint8(config.LowPassFilterConfig) << 3)
	}
	var sensitivity uint8 = 0
	switch config.SensitivityLevel {
	case "250dps":
		sensitivity = 0
	case "500dps":
		sensitivity = 1
	case "1000dps":
		sensitivity = 2
	case "2000dps":
		sensitivity = 3
	}

	gyroConfig1 = gyroConfig1 | (uint8(sensitivity) << 1)
	err := dev.writeRegister(GYRO_CONFIG_1, gyroConfig1, gyroConfig2)
	cnfg, _ := dev.readRegister(GYRO_CONFIG_1, 2)
	dev.setGyroOffset(XG_OFFS_USRH, config.Offsets.X)
	dev.setGyroOffset(YG_OFFS_USRH, config.Offsets.Y)
	dev.setGyroOffset(ZG_OFFS_USRH, config.Offsets.Z)
	fmt.Println("Gyro Config: ", cnfg)
	return err
}

func (dev *memsICM20948) processGyroscopeData(data []uint8) (types.XYZ, error) {
	gyroConfig, _ := dev.GetGyro().GetConfig().(types.GyroscopeConfig)
	scale := gyroFullScale[gyroConfig.SensitivityLevel]
	x := float64(utils.TowsComplementUint8ToInt16(data[0], data[1])) / scale
	y := float64(utils.TowsComplementUint8ToInt16(data[2], data[3])) / scale
	z := float64(utils.TowsComplementUint8ToInt16(data[4], data[5])) / scale
	return types.XYZ{
		X: x,
		Y: y,
		Z: z,
	}, nil
}

func (dev *memsICM20948) setGyroOffset(address uint16, offset int16) {
	var h uint8 = uint8(uint16(offset) >> 8)
	var l uint8 = uint8(uint16(offset) & 0xFF)
	dev.writeRegister(address, h, l)
}
