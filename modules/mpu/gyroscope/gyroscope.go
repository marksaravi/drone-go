package gyroscope

// Gyroscope is an interface to gyroscope mems operations
type Gyroscope interface {
	//SetFullScaleRange is to set gyroscope range
	SetFullScaleRange(byte) error

	//GetFullScaleRange is to get gyroscope range
	GetFullScaleRange() (byte, error)
}
