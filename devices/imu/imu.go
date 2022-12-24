package imu

import "github.com/marksaravi/drone-go/types"

type IMUMems6DOF interface {
	Setup()
	ReadIMUData() (types.IMUMems6DOFData, error)
}

type imuDevice struct {
	dev IMUMems6DOF
}

func NewIMU(dev IMUMems6DOF) *imuDevice {
	return &imuDevice{
		dev: dev,
	}
}

func (imu *imuDevice) Setup() {
	imu.dev.Setup()
}

func (imu *imuDevice) Read() (types.Orientations, bool) {
	data, err := imu.dev.ReadIMUData()
	if err != nil {
		return types.Orientations{}, false
	}
	return imu.calcRotations(data), true
}

func (imu *imuDevice) calcRotations(types.IMUMems6DOFData) types.Orientations {
	return types.Orientations{}
}
