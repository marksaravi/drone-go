package imu

import (
	"math"
	"time"

	"github.com/marksaravi/drone-go/types"
)

const MIN_TIME_BETWEEN_READS = time.Millisecond * 50

type IMUMems6DOF interface {
	Read() (types.IMUMems6DOFRawData, error)
}

type imuDevice struct {
	dev               IMUMems6DOF
	rotations         types.Rotations
	lastReadTime      time.Time
	currReadTime      time.Time
	filterCoefficient float64
}

func NewIMU(dev IMUMems6DOF, configs types.IMUConfigs) *imuDevice {
	return &imuDevice{
		dev: dev,
		rotations: types.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		filterCoefficient: configs.FilterCoefficient,
	}
}

// Read returns Roll, Pitch and Yaw.
func (imu *imuDevice) Read() (types.Rotations, error) {
	imu.currReadTime = time.Now()
	data, err := imu.dev.Read()
	if err != nil {
		return types.Rotations{}, err
	}
	imu.calcRotations(data)
	imu.lastReadTime = imu.currReadTime
	return imu.rotations, nil
}

func (imu *imuDevice) calcRotations(memsData types.IMUMems6DOFRawData) {
	acc := calcaAcelerometerRotations(memsData.Accelerometer)
	dt := imu.currReadTime.Sub(imu.lastReadTime)
	gyro := calcGyroscopeRotations(memsData.Gyroscope, dt, imu.rotations)
	imu.rotations = types.Rotations{
		Roll:  complimentaryFilter(gyro.Roll, acc.Roll, imu.filterCoefficient),
		Pitch: complimentaryFilter(gyro.Pitch, acc.Pitch, imu.filterCoefficient),
		Yaw:   gyro.Yaw,
	}
}

func complimentaryFilter(gyro float64, accelerometer float64, complimentaryFilterCoefficient float64) float64 {
	return (1-complimentaryFilterCoefficient)*gyro + complimentaryFilterCoefficient*accelerometer
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

func calcGyroscopeRotations(gyroData types.DXYZ, dt time.Duration, prevRotations types.Rotations) types.Rotations {
	if dt > MIN_TIME_BETWEEN_READS {
		return prevRotations
	}
	roll := prevRotations.Roll + gyroData.DX*dt.Seconds()
	pitch := prevRotations.Pitch + gyroData.DY*dt.Seconds()
	yaw := prevRotations.Yaw + gyroData.DY*dt.Seconds()
	return types.Rotations{
		Roll:  math.Mod(roll, 360),
		Pitch: math.Mod(pitch, 360),
		Yaw:   math.Mod(yaw, 360),
	}
}
