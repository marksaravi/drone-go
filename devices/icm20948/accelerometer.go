package icm20948

import (
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
	var accConfig1 uint8 = 0b00000000
	var accConfig2 uint8 = uint8(config.Averaging)
	if config.LowPassFilterEnabled {
		accConfig1 = 0b00000001 | (uint8(config.LowPassFilterConfig) << 3)
	}
	accConfig1 = accConfig1 | (uint8(config.SensitivityLevel) << 2)
	err := dev.writeRegister(ACCEL_CONFIG, accConfig1, accConfig2)
	time.Sleep(time.Millisecond * 100)
	return err
}

func (dev *MemsICM20948) getAccConfig() (AccelerometerConfig, error) {
	data, err := dev.readRegister(ACCEL_CONFIG, 2)
	config := AccelerometerConfig{
		SensitivityLevel: int((data[0] >> 1) & 0b00000011),
	}
	dev.GetAcc().SetConfig(config)
	return config, err
}

func (dev *MemsICM20948) processAccelerometerData(data []byte) (types.XYZ, error) {
	accConfig, _ := dev.GetAcc().GetConfig().(AccelerometerConfig)
	accSens := accelerometerSensitivity[accConfig.SensitivityLevel]
	x := float64(utils.TowsComplementBytesToInt(data[0], data[1]))/accSens - accConfig.Offsets.X
	y := float64(utils.TowsComplementBytesToInt(data[2], data[3]))/accSens - accConfig.Offsets.Y
	z := float64(utils.TowsComplementBytesToInt(data[4], data[5]))/accSens - accConfig.Offsets.Z
	return types.XYZ{
		X: x,
		Y: y,
		Z: z,
	}, nil
}
