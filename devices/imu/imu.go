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
	DataPerSecond                               int     `json:"data-per-second"`
	AccelerometerComplimentaryFilterCoefficient float64 `json:"accelerometer-complimentary-filter-coefficient"`
	RotationsComplimentaryFilterCoefficient     float64 `json:"rotations-complimentary-filter-coefficient"`
}

// Rotations (Roll, Pitch, Yaw)
type Rotations struct {
	Roll, Pitch, Yaw float64
}

type ImuData struct {
	Accelerometer Rotations
	Gyroscope     Rotations
	Rotations     Rotations
	Error         error
}

type imuDevice struct {
	configs                           Configs
	dev                               IMUMems6DOF
	rotations                         Rotations
	accRotations                      Rotations
	gyroRotations                     Rotations
	dRoll                             float64
	dPitch                            float64
	dYaw                              float64
	lastReadTime                      time.Time
	currReadTime                      time.Time
	accComplimentaryFilterCoefficient float64
	rotComplimentaryFilterCoefficient float64
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
		dRoll:                             0,
		dPitch:                            0,
		dYaw:                              0,
		accComplimentaryFilterCoefficient: configs.AccelerometerComplimentaryFilterCoefficient,
		rotComplimentaryFilterCoefficient: configs.RotationsComplimentaryFilterCoefficient,
	}
}

func (imuDev *imuDevice) Reset() {
	imuDev.currReadTime = time.Now()
	imuDev.lastReadTime = imuDev.currReadTime
}

func (imuDev *imuDevice) Read() (Rotations, Rotations, Rotations, error) {
	imuDev.currReadTime = time.Now()
	data, err := imuDev.dev.Read()
	if err != nil {
		return imuDev.rotations, imuDev.accRotations, imuDev.gyroRotations, err
	}
	imuDev.calcAllRotations(data)
	imuDev.lastReadTime = imuDev.currReadTime
	return imuDev.rotations, imuDev.accRotations, imuDev.gyroRotations, nil
}

func (imuDev *imuDevice) calcAllRotations(memsData mems.Mems6DOFData) {
	imuDev.calcaAccelerometerRotations(memsData.Accelerometer)
	imuDev.calcGyroscopeRotations(memsData.Gyroscope)
	imuDev.calcRotations()
}

func (imuDev *imuDevice) calcaAccelerometerRotations(data mems.XYZ) {
	pitch := 180 * math.Atan2(data.X, math.Sqrt(data.Y*data.Y+data.Z*data.Z)) / math.Pi
	roll := 180 * math.Atan2(data.Y, math.Sqrt(data.X*data.X+data.Z*data.Z)) / math.Pi
	imuDev.accRotations = Rotations{
		Roll:  complimentaryFilter(roll, imuDev.accRotations.Roll, imuDev.accComplimentaryFilterCoefficient),
		Pitch: complimentaryFilter(pitch, imuDev.accRotations.Pitch, imuDev.accComplimentaryFilterCoefficient),
		Yaw:   0,
	}
}

func (imuDev *imuDevice) calcGyroscopeRotations(dxyz mems.DXYZ) {
	dt := imuDev.currReadTime.Sub(imuDev.lastReadTime)
	if dt < MIN_TIME_BETWEEN_READS {
		return
	}

	imuDev.dRoll = dxyz.DX * dt.Seconds()
	imuDev.dPitch = dxyz.DY * dt.Seconds()
	imuDev.dYaw = dxyz.DZ * dt.Seconds()

	imuDev.gyroRotations.Roll += imuDev.dRoll
	imuDev.gyroRotations.Pitch += imuDev.dPitch
	imuDev.gyroRotations.Yaw += imuDev.dYaw
}

func (imuDev *imuDevice) calcRotations() {
	roll := imuDev.rotations.Roll + imuDev.dRoll
	pitch := imuDev.rotations.Pitch + imuDev.dPitch
	yaw := imuDev.rotations.Yaw + imuDev.dYaw

	imuDev.rotations = Rotations{
		Roll:  complimentaryFilter(roll, imuDev.accRotations.Roll, imuDev.rotComplimentaryFilterCoefficient),
		Pitch: complimentaryFilter(pitch, imuDev.accRotations.Pitch, imuDev.rotComplimentaryFilterCoefficient),
		Yaw:   yaw,
	}
}

func complimentaryFilter(value float64, preValue float64, complimentaryFilterCoefficient float64) float64 {
	v := (1-complimentaryFilterCoefficient)*value + complimentaryFilterCoefficient*preValue
	return v
}
