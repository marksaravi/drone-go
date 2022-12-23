package imu

type inertialDevice interface {
	Setup()
	ReadIMUData() ([]byte, error)
}

type Rotations struct {
	Roll, Pitch, Yaw float64
}

type imuDevice struct {
	dev inertialDevice
}

func NewIMU(dev inertialDevice) *imuDevice {
	return &imuDevice{
		dev: dev,
	}
}

func (imu *imuDevice) Setup() {
	imu.dev.Setup()
}

func (imu *imuDevice) ReadInertialDevice() (Rotations, bool) {
	data, err := imu.dev.ReadIMUData()
	if err != nil {
		return Rotations{}, false
	}
	return imu.calcRotations(data), true
}

func (imu *imuDevice) calcRotations(data []byte) Rotations {
	return Rotations{}
}
