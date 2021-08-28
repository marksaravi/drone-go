package icm20948

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/models"
)

// InitGyroscope initialise the Gyroscope
func (dev *memsICM20948) InitGyroscope() error {
	var gyroConfig1 uint8 = 0b00000000
	var gyroConfig2 uint8 = uint8(dev.gyroConfig.averaging)

	if dev.gyroConfig.lowPassFilterEnabled {
		gyroConfig1 = 0b00000001 | (uint8(dev.gyroConfig.lowPassFilterConfig) << 3)
	}
	var sensitivity uint8 = 0
	switch dev.gyroConfig.sensitivityLevel {
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
	dev.setGyroOffset(XG_OFFS_USRH, dev.gyroConfig.offsetX)
	dev.setGyroOffset(YG_OFFS_USRH, dev.gyroConfig.offsetY)
	dev.setGyroOffset(ZG_OFFS_USRH, dev.gyroConfig.offsetZ)
	fmt.Println("Gyro Config: ", cnfg)
	return err
}

func (dev *memsICM20948) processGyroscopeData(data []uint8) (models.XYZ, error) {
	scale := gyroFullScale[dev.gyroConfig.sensitivityLevel]
	x := float64(towsComplementUint8ToInt16(data[0], data[1])) / scale
	y := float64(towsComplementUint8ToInt16(data[2], data[3])) / scale
	z := float64(towsComplementUint8ToInt16(data[4], data[5])) / scale
	return models.XYZ{
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
