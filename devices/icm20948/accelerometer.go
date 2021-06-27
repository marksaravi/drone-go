package icm20948

import (
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

// GetAcc get accelerometer data
func (dev *MemsICM20948) GetAcc() *types.Sensor {
	return &(dev.acc)
}

// InitAccelerometer initialise the Accelerometer
func (dev *MemsICM20948) InitAccelerometer() error {
	config, ok := dev.GetAcc().GetConfig().(AccelerometerConfig)
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
	time.Sleep(time.Millisecond * 100)
	return err
}

func (dev *MemsICM20948) processAccelerometerData(data []byte) (types.XYZ, error) {
	config, _ := dev.GetAcc().GetConfig().(AccelerometerConfig)
	accSens := accelerometerSensitivity[config.SensitivityLevel]
	xRaw := utils.TowsComplementBytesToInt(data[0], data[1])
	yRaw := utils.TowsComplementBytesToInt(data[2], data[3])
	zRaw := utils.TowsComplementBytesToInt(data[4], data[5])
	x := float64(xRaw)/accSens - config.Offsets.X
	y := float64(yRaw)/accSens - config.Offsets.Y
	z := float64(zRaw)/accSens - config.Offsets.Z
	fmt.Println(zRaw, z, accSens)
	return types.XYZ{
		X: x,
		Y: y,
		Z: z,
	}, nil
}
