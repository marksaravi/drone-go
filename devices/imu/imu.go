package imu

import (
	"math"
	"time"

	"github.com/marksaravi/drone-go/hardware/mems"
)

const MIN_TIME_BETWEEN_READS = time.Nanosecond

type IMUMems6DOF interface {
	Read() (mems.Mems6DOFData, error)
}

type Configs struct {
	AccelerometerComplimentaryFilterCoefficient float64 `yaml:"acc-complimentary_filter_coefficient"`
}

// Rotations (Roll, Pitch, Yaw)
type Rotations struct {
	Roll, Pitch, Yaw float64
}
type imuDevice struct {
	configs                           Configs
	dev                               IMUMems6DOF
	rotations                         Rotations
	accRotations                      Rotations
	gyroRotations                     Rotations
	lastReadTime                      time.Time
	currReadTime                      time.Time
	accComplimentaryFilterCoefficient float64
}

func NewIMU(dev IMUMems6DOF, configs Configs) *imuDevice {
	return &imuDevice{
		configs: configs,
		dev:     dev,
		rotations: Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		accRotations: Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		gyroRotations: Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		accComplimentaryFilterCoefficient: configs.AccelerometerComplimentaryFilterCoefficient,
	}
}

// Read returns Roll, Pitch and Yaw.
func (imu *imuDevice) Read() (Rotations, Rotations, Rotations, error) {
	imu.currReadTime = time.Now()
	data, err := imu.dev.Read()
	if err != nil {
		return imu.rotations, imu.accRotations, imu.gyroRotations, err
	}
	imu.calcRotations(data)
	imu.lastReadTime = imu.currReadTime
	return imu.rotations, imu.accRotations, imu.gyroRotations, nil
}

func (imu *imuDevice) calcRotations(memsData mems.Mems6DOFData) {
	imu.calcaAcelerometerRotations(memsData.Accelerometer)
	imu.calcGyroscopeRotations(memsData.Gyroscope)
}

func (imu *imuDevice) calcaAcelerometerRotations(data mems.XYZ) {
	pitch := 180 * math.Atan2(data.X, math.Sqrt(data.Y*data.Y+data.Z*data.Z)) / math.Pi
	roll := 180 * math.Atan2(data.Y, math.Sqrt(data.X*data.X+data.Z*data.Z)) / math.Pi
	imu.accRotations = Rotations{
		Roll:  complimentaryFilter(roll, imu.accRotations.Roll, imu.accComplimentaryFilterCoefficient),
		Pitch: complimentaryFilter(pitch, imu.accRotations.Pitch, imu.accComplimentaryFilterCoefficient),
		Yaw:   0,
	}
}

func (imu *imuDevice) calcGyroscopeRotations(dxyz mems.DXYZ) {
	dt := imu.currReadTime.Sub(imu.lastReadTime)
	if dt > MIN_TIME_BETWEEN_READS {
		return
	}
	dRoll := dxyz.DX*dt.Seconds()
	dPitch := dxyz.DY*dt.Seconds()
	dYaw := dxyz.DZ*dt.Seconds()
	imu.gyroRotations.Roll += dRoll
	imu.gyroRotations.Pitch += dPitch
	imu.gyroRotations.Yaw += dYaw
}

func complimentaryFilter(value float64, preValue float64, complimentaryFilterCoefficient float64) float64 {
	v := (1-complimentaryFilterCoefficient)*value + complimentaryFilterCoefficient*preValue
	return math.Round(v*10)/10
}