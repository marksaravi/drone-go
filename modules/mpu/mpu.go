package mpu

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	ResetToDefault() error
	WhoAmI() (string, byte, error)
}
