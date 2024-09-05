package imu

import (
	"math"
	"time"

	"github.com/marksaravi/drone-go/hardware/mems"
)

const MIN_TIME_BETWEEN_READS = time.Second / 10000

type IMUMems6DOF interface {
	Read() (mems.Mems6DOFData, error)
}

type Directions struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Configs struct {
	DataPerSecond                           int        `json:"data-per-second"`
	RotationsComplimentaryFilterCoefficient float64    `json:"rotations-complimentary-filter-coefficient"`
	AccDirections                           Directions `json:"acc-directions"`
	GyroscopDirections                      Directions `json:"gyro-directions"`
}

// Rotations (Roll, Pitch, Yaw)
type Rotations struct {
	Roll, Pitch, Yaw float64
	Time             time.Time
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
	firstRead                         bool
	lastReadTime                      time.Time
	currReadTime                      time.Time
	rotComplimentaryFilterCoefficient float64
	accDirX                           float64
	accDirY                           float64
	accDirZ                           float64
	gyroDirX                          float64
	gyroDirY                          float64
	gyroDirZ                          float64
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
		rotComplimentaryFilterCoefficient: configs.RotationsComplimentaryFilterCoefficient,
		firstRead:                         true,
		accDirX:                           configs.AccDirections.X,
		accDirY:                           configs.AccDirections.Y,
		accDirZ:                           configs.AccDirections.Z,
		gyroDirX:                          configs.GyroscopDirections.X,
		gyroDirY:                          configs.GyroscopDirections.Y,
		gyroDirZ:                          configs.GyroscopDirections.Z,
	}
}

func (imuDev *imuDevice) Read() (Rotations, error) {
	r, _, _, err := imuDev.ReadAll()
	return r, err
}

func (imuDev *imuDevice) ReadAll() (Rotations, Rotations, Rotations, error) {
	data, err := imuDev.dev.Read()
	imuDev.currReadTime = data.Time
	if err != nil {
		return imuDev.rotations, imuDev.accRotations, imuDev.gyroRotations, err
	}
	imuDev.calcAllRotations(data)
	imuDev.lastReadTime = imuDev.currReadTime
	imuDev.rotations.Time = imuDev.currReadTime
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
		Roll:  roll * imuDev.accDirX,
		Pitch: pitch * imuDev.accDirY,
		Yaw:   0,
	}
}

func truncate(x, max float64) float64 {
	if math.Abs(x) < max {
		return x
	}
	if x < 0 {
		return x + max
	}
	return x - max
}

func (imuDev *imuDevice) calcGyroscopeRotations(dxyz mems.DXYZ) {
	dt := imuDev.currReadTime.Sub(imuDev.lastReadTime)
	if dt < MIN_TIME_BETWEEN_READS || imuDev.firstRead {
		imuDev.firstRead = false
		return
	}

	imuDev.dRoll = dxyz.DX * dt.Seconds()
	imuDev.dPitch = dxyz.DY * dt.Seconds()
	imuDev.dYaw = dxyz.DZ * dt.Seconds()

	imuDev.gyroRotations.Roll = truncate(imuDev.gyroRotations.Roll+imuDev.dRoll*imuDev.gyroDirX, 360)
	imuDev.gyroRotations.Pitch = truncate(imuDev.gyroRotations.Pitch+imuDev.dPitch*imuDev.gyroDirY, 360)
	imuDev.gyroRotations.Yaw = truncate(imuDev.gyroRotations.Yaw+imuDev.dYaw*imuDev.gyroDirZ, 360)
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
