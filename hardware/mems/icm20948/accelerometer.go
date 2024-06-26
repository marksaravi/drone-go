package icm20948

import (
	"time"

	"github.com/marksaravi/drone-go/models"
)

// initAccelerometer initialise the Accelerometer
func (dev *memsICM20948) initAccelerometer() error {
	var sensitivity uint8 = 0
	switch dev.accConfig.sensitivityLevel {
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
	var accConfig2 uint8 = uint8(dev.accConfig.averaging)
	if dev.accConfig.lowPassFilterEnabled {
		accConfig1 = accConfig1 | 0b00000001
		accConfig1 = accConfig1 | (uint8(dev.accConfig.lowPassFilterConfig) << 3)
	}
	err := dev.writeRegister(ACCEL_CONFIG, accConfig1, accConfig2)
	dev.setAccOffset(XA_OFFS_H, dev.accConfig.offsetX)
	dev.setAccOffset(YA_OFFS_H, dev.accConfig.offsetY)
	dev.setAccOffset(ZA_OFFS_H, dev.accConfig.offsetZ)
	time.Sleep(time.Millisecond * 100)
	return err
}

func accOffsetToHighandLowBytes(offset uint16) (h, l uint8) {
	shiftedOffset := offset << 1
	h = uint8((shiftedOffset >> 8) & 0xFF)
	l = uint8(shiftedOffset & 0xFF)
	return
}

func (dev *memsICM20948) setAccOffset(addressH uint16, offset uint16) {
	h, l := accOffsetToHighandLowBytes(offset)
	dev.writeRegister(addressH, h, l)
}

func (dev *memsICM20948) processAccelerometerData(data []byte) (models.XYZ, error) {
	accSens := accelerometerSensitivity[dev.accConfig.sensitivityLevel]
	xRaw := towsComplementUint8ToInt16(data[0], data[1])
	yRaw := towsComplementUint8ToInt16(data[2], data[3])
	zRaw := towsComplementUint8ToInt16(data[4], data[5])
	x := float64(xRaw) / accSens
	y := float64(yRaw) / accSens
	z := float64(zRaw) / accSens

	return models.XYZ{
		X: x,
		Y: y,
		Z: z,
	}, nil
}
