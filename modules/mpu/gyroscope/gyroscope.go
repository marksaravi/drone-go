package gyroscope

// Gyroscope is an interface to gyroscope mems operations
type Gyroscope interface {
	//SetRegister is to set gyroscope range
	SetRegister(address, bank, data byte) error

	//GetRegister is to get gyroscope range
	GetRegister(address, bank byte) ([]byte, error)
}
