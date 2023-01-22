package imu

import (
	"math"
	"time"

	"github.com/marksaravi/drone-go/types"
	"github.com/marksaravi/drone-go/utils"
)

const MIN_TIME_BETWEEN_READS = time.Millisecond * 50

type IMUMems6DOF interface {
	Read() (types.IMUMems6DOFRawData, error)
}

type Configs struct {
	DataPerSecond     int     `yaml:"data_per_second"`
	FilterCoefficient float64 `yaml:"filter_coefficient"`
}

// Rotations (Roll, Pitch, Yaw)
type Rotations struct {
	Roll, Pitch, Yaw float64
}
type imuDevice struct {
	configs           Configs
	dev               IMUMems6DOF
	rotations         Rotations
	lastReadTime      time.Time
	currReadTime      time.Time
	filterCoefficient float64
}

func NewIMU(dev IMUMems6DOF) *imuDevice {
	configs := readConfigs()
	return &imuDevice{
		configs: configs,
		dev:     dev,
		rotations: Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		filterCoefficient: configs.FilterCoefficient,
	}
}

func readConfigs() Configs {
	var configs struct {
		Imu Configs `yaml:"imu"`
	}
	utils.ReadConfigs(&configs)
	return configs.Imu
}

// Read returns Roll, Pitch and Yaw.
func (imu *imuDevice) Read() (Rotations, error) {
	imu.currReadTime = time.Now()
	data, err := imu.dev.Read()
	if err != nil {
		return Rotations{}, err
	}
	imu.calcRotations(data)
	imu.lastReadTime = imu.currReadTime
	return imu.rotations, nil
}

func (imu *imuDevice) calcRotations(memsData types.IMUMems6DOFRawData) {
	acc := calcaAcelerometerRotations(memsData.Accelerometer)
	dt := imu.currReadTime.Sub(imu.lastReadTime)
	gyro := calcGyroscopeRotations(memsData.Gyroscope, dt, imu.rotations)
	imu.rotations = Rotations{
		Roll:  complimentaryFilter(gyro.Roll, acc.Roll, imu.filterCoefficient),
		Pitch: complimentaryFilter(gyro.Pitch, acc.Pitch, imu.filterCoefficient),
		Yaw:   gyro.Yaw,
	}
}

func complimentaryFilter(gyro float64, accelerometer float64, complimentaryFilterCoefficient float64) float64 {
	return (1-complimentaryFilterCoefficient)*gyro + complimentaryFilterCoefficient*accelerometer
}

func calcaAcelerometerRotations(data types.XYZ) Rotations {
	yrot := 180 * math.Atan2(data.X, math.Sqrt(data.Y*data.Y+data.Z*data.Z)) / math.Pi
	xrot := 180 * math.Atan2(data.Y, math.Sqrt(data.X*data.X+data.Z*data.Z)) / math.Pi
	return Rotations{
		Roll:  xrot,
		Pitch: yrot,
		Yaw:   0,
	}
}

func calcGyroscopeRotations(gyroData types.DXYZ, dt time.Duration, prevRotations Rotations) Rotations {
	if dt > MIN_TIME_BETWEEN_READS {
		return prevRotations
	}
	roll := prevRotations.Roll + gyroData.DX*dt.Seconds()
	pitch := prevRotations.Pitch + gyroData.DY*dt.Seconds()
	yaw := prevRotations.Yaw + gyroData.DZ*dt.Seconds()
	return Rotations{
		Roll:  math.Mod(roll, 360),
		Pitch: math.Mod(pitch, 360),
		Yaw:   math.Mod(yaw, 360),
	}
}
