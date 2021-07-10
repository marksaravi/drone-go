package hardware

import "github.com/MarkSaravi/drone-go/hardware/icm20948"

type HardwareConfig struct {
	ICM20948 icm20948.Icm20948Config `yaml:"icm20948"`
}
