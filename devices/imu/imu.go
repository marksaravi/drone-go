package imu

import (
	"log"
	"math"

	"github.com/marksaravi/drone-go/types"
)

type IMUMems6DOF interface {
	ReadIMUData() (types.IMUMems6DOFRawData, error)
}

type imuDevice struct {
	dev IMUMems6DOF

	rotations types.Rotations

	compFilteCoefficient float64
}

func NewIMU(dev IMUMems6DOF) *imuDevice {
	return &imuDevice{
		dev: dev,
		rotations: types.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		compFilteCoefficient: 0.001,
	}
}

func (imu *imuDevice) Read() (types.Rotations, bool) {
	data, err := imu.dev.ReadIMUData()
	if err != nil {
		return types.Rotations{}, false
	}
	log.Print(data)
	return imu.calcRotations(data), true
}

func (imu *imuDevice) calcRotations(memsData types.IMUMems6DOFRawData) types.Rotations {
	acc := calcaAcelerometerRotations(memsData.Accelerometer)
	gyro := calcGyroscopeRotations(memsData.Gyroscope, imu.rotations)
	return types.Rotations{
		Roll:  complimentaryFilter(gyro.Roll, acc.Roll, imu.compFilteCoefficient),
		Pitch: complimentaryFilter(gyro.Pitch, acc.Pitch, imu.compFilteCoefficient),
		Yaw:   gyro.Yaw,
	}
}

func complimentaryFilter(gyroValue float64, accelerometerValue float64, complimentaryFilterCoefficient float64) float64 {
	// v1 := (1 - complimentaryFilterCoefficient) * gyroValue
	// v2 := complimentaryFilterCoefficient * accelerometerValue
	// return v1 + v2
	return accelerometerValue
}

func calcaAcelerometerRotations(data types.XYZ) types.Rotations {
	yrot := 180 * math.Atan2(data.X, math.Sqrt(data.Y*data.Y+data.Z*data.Z)) / math.Pi
	xrot := 180 * math.Atan2(data.Y, math.Sqrt(data.X*data.X+data.Z*data.Z)) / math.Pi
	return types.Rotations{
		Roll:  xrot,
		Pitch: yrot,
		Yaw:   0,
	}
}

func calcGyroscopeRotations(dGyro types.XYZDt, prevRotations types.Rotations) types.Rotations {
	// return types.Rotations{
	// 	Roll:  math.Mod(prevRotations.Roll+dGyro.DX, 360),
	// 	Pitch: math.Mod(prevRotations.Roll+dGyro.DY, 360),
	// 	Yaw:   math.Mod(prevRotations.Roll+dGyro.DZ, 360),
	// }
	return types.Rotations{
		Roll:  0,
		Pitch: 0,
		Yaw:   0,
	}
}
