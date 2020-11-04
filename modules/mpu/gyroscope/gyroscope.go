package gyroscope

// Gyroscope is an interface to gyroscope mems operations
type Gyroscope interface {
	//SetFullScaleRange is to set gyroscope range
	SetFullScaleRange()

	//GetFullScaleRange is to get gyroscope range
	GetFullScaleRange()
}
