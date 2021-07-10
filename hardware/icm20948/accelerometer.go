package icm20948

import (
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

// GetAcc get accelerometer data
func (dev *memsICM20948) GetAcc() *types.Sensor {
	return &(dev.acc)
}

// InitAccelerometer initialise the Accelerometer
func (dev *memsICM20948) InitAccelerometer() error {
	config, ok := dev.GetAcc().GetConfig().(accelerometerConfig)
	if !ok {
		log.Fatal("Accelerometer config mismatch")
	}
	var sensitivity uint8 = 0
	switch config.SensitivityLevel {
	case "2g":
		sensitivity = 0
	case "4g":
		sensitivity = 1
	case "8g":
		sensitivity = 2
	case "16g":
		sensitivity = 3
	}
	var accConfig1 uint8 = 0b00000000 | (uint8(sensitivity) << 1)
	var accConfig2 uint8 = uint8(config.Averaging)
	if config.LowPassFilterEnabled {
		accConfig1 = accConfig1 | 0b00000001
		accConfig1 = accConfig1 | (uint8(config.LowPassFilterConfig) << 3)
	}
	err := dev.writeRegister(ACCEL_CONFIG, accConfig1, accConfig2)
	dev.setAccOffset(XA_OFFS_H, config.Offsets.X)
	dev.setAccOffset(YA_OFFS_H, config.Offsets.Y)
	dev.setAccOffset(ZA_OFFS_H, config.Offsets.Z)
	time.Sleep(time.Millisecond * 100)
	return err
}

func (dev *memsICM20948) processAccelerometerData(data []byte) (types.XYZ, error) {
	config, _ := dev.GetAcc().GetConfig().(accelerometerConfig)
	accSens := accelerometerSensitivity[config.SensitivityLevel]
	xRaw := utils.TowsComplementUint8ToInt16(data[0], data[1])
	yRaw := utils.TowsComplementUint8ToInt16(data[2], data[3])
	zRaw := utils.TowsComplementUint8ToInt16(data[4], data[5])
	x := float64(xRaw) / accSens
	y := float64(yRaw) / accSens
	z := float64(zRaw) / accSens

	return types.XYZ{
		X: x,
		Y: y,
		Z: z,
	}, nil
}

func (dev *memsICM20948) setAccOffset(address uint16, offset int16) {
	offsets, _ := dev.readRegister(address, 2)
	var h uint8 = uint8((uint16(offset) >> 7) & 0xFF)
	var l uint8 = uint8((uint16(offset)<<1)&0xFF) | (offsets[1] & 0x01)
	dev.writeRegister(address, h, l)
}
