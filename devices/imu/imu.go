package imu

import (
	"math"
	"time"

	"github.com/marksaravi/drone-go/hardware/mems"
)

const MIN_TIME_BETWEEN_READS = time.Millisecond * 50

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
	// dt := imu.currReadTime.Sub(imu.lastReadTime)
	// var rotations Rotations
	// rotations, imu.gyroRotations = calcGyroscopeRotations(memsData.Gyroscope, dt, imu.rotations, imu.gyroRotations)

	// imu.rotations = Rotations{
	// 	Roll:  complimentaryFilter(rotations.Roll, imu.accRotations.Roll, imu.complimentaryFilterCoefficient),
	// 	Pitch: complimentaryFilter(rotations.Pitch, imu.accRotations.Pitch, imu.complimentaryFilterCoefficient),
	// 	Yaw:   rotations.Yaw,
	// }
}

func (imu *imuDevice) calcaAcelerometerRotations(data mems.XYZ) {
	// fmt.Println(imu.accComplimentaryFilterCoefficient)
	pitch := 180 * math.Atan2(data.X, math.Sqrt(data.Y*data.Y+data.Z*data.Z)) / math.Pi
	roll := 180 * math.Atan2(data.Y, math.Sqrt(data.X*data.X+data.Z*data.Z)) / math.Pi
	imu.accRotations = Rotations{
		Roll:  complimentaryFilter(roll, imu.accRotations.Roll, imu.accComplimentaryFilterCoefficient),
		Pitch: complimentaryFilter(pitch, imu.accRotations.Pitch, imu.accComplimentaryFilterCoefficient),
		Yaw:   0,
	}
}

func calcGyroscopeRotations(gyroData mems.DXYZ, dt time.Duration, prevRotations Rotations, gyroRotations Rotations) (Rotations, Rotations) {
	if dt > MIN_TIME_BETWEEN_READS {
		return prevRotations, gyroRotations
	}
	roll := prevRotations.Roll + gyroData.DX*dt.Seconds()
	pitch := prevRotations.Pitch + gyroData.DY*dt.Seconds()
	yaw := prevRotations.Yaw + gyroData.DZ*dt.Seconds()
	gyroRoll := gyroRotations.Roll + gyroData.DX*dt.Seconds()
	gyroPitch := gyroRotations.Pitch + gyroData.DY*dt.Seconds()
	gyroYaw := gyroRotations.Yaw + gyroData.DZ*dt.Seconds()

	return Rotations{
			Roll:  math.Mod(roll, 360),
			Pitch: math.Mod(pitch, 360),
			Yaw:   math.Mod(yaw, 360),
		}, Rotations{
			Roll:  math.Mod(gyroRoll, 360),
			Pitch: math.Mod(gyroPitch, 360),
			Yaw:   math.Mod(gyroYaw, 360),
		}
}

func complimentaryFilter(value float64, preValue float64, complimentaryFilterCoefficient float64) float64 {
	v := (1-complimentaryFilterCoefficient)*value + complimentaryFilterCoefficient*preValue
	return math.Round(v*10)/10
}