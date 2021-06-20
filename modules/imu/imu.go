package imu

import (
	"fmt"
	"math"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

type ImuModule struct {
	dev                      types.ImuDevice
	rotations                types.Rotations
	gyro                     types.Rotations
	startTime                time.Time
	readTime                 time.Time
	readingInterval          time.Duration
	lowPassFilterCoefficient float64
	readingData              types.ImuReadingQualities
}

func (imu *ImuModule) ResetReadingTimes() {
	imu.startTime = time.Now()
	imu.readTime = imu.startTime
	imu.rotations = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
	imu.gyro = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
}

func (imu *ImuModule) updateReadingData(diff time.Duration, err bool) {
	imu.readingData.Total++
	if diff >= imu.readingData.BadIntervalThereshold {
		imu.readingData.BadInterval++
		if diff > imu.readingData.MaxBadInterval {
			imu.readingData.MaxBadInterval = diff
		}
	}
	if err {
		imu.readingData.BadData++
	}
}

func (imu *ImuModule) readSensors() (acc types.SensorData, gyro types.SensorData, mag types.SensorData, err error) {
	acc, gyro, mag, err = imu.dev.ReadSensors()
	return
}

func NewIMU(imuMems types.ImuDevice, config types.FlightConfig) ImuModule {
	readingInterval := time.Duration(int64(time.Second) / int64(config.ImuDataPerSecond))
	badIntervalThereshold := readingInterval + readingInterval/20
	fmt.Println(readingInterval, badIntervalThereshold)
	return ImuModule{
		dev:                      imuMems,
		readTime:                 time.Time{},
		readingInterval:          readingInterval,
		lowPassFilterCoefficient: config.LowPassFilterCoefficient,
		readingData: types.ImuReadingQualities{
			Total:                 0,
			BadInterval:           0,
			MaxBadInterval:        0,
			BadData:               0,
			BadIntervalThereshold: badIntervalThereshold,
		},
	}
}

func (imu ImuModule) Close() {
	imu.dev.Close()
}

func (imu *ImuModule) GetRotations() (bool, types.ImuRotations, error) {
	now := time.Now()
	diff := now.Sub(imu.readTime)
	if diff < imu.readingInterval {
		return false, types.ImuRotations{}, nil
	}
	imu.readTime = now
	acc, gyro, _, imuError := imu.readSensors()
	imu.updateReadingData(diff, imuError != nil)
	dg := GyroChanges(gyro, diff.Nanoseconds())
	imu.gyro = GyroRotations(dg, imu.gyro)
	accRotations := AccelerometerRotations(acc.Data)
	prevRotations := imu.rotations
	imu.rotations = CalcRotations(
		prevRotations,
		accRotations,
		dg,
		imu.lowPassFilterCoefficient,
	)

	return true, types.ImuRotations{
		Accelerometer: accRotations,
		Gyroscope:     imu.gyro,
		Rotations:     imu.rotations,
		ReadTime:      imu.readTime.UnixNano() - imu.startTime.UnixNano(),
		ReadInterval:  diff.Nanoseconds(),
	}, imuError
}

func (imu *ImuModule) GetReadingQualities() types.ImuReadingQualities {
	return imu.readingData
}

func GyroChanges(gyro types.SensorData, timeInterval int64) types.RotationsChanges {
	dt := goDurToDt(timeInterval)
	return types.RotationsChanges{
		DRoll:  gyro.Data.X * dt,
		DPitch: gyro.Data.Y * dt,
		DYaw:   gyro.Data.Z * dt,
	}
}

func goDurToDt(d int64) float64 {
	return float64(d) / 1e9
}

func GyroRotations(dg types.RotationsChanges, gyroRotations types.Rotations) types.Rotations {
	return types.Rotations{
		Roll:  math.Mod(gyroRotations.Roll+dg.DRoll, 360),
		Pitch: math.Mod(gyroRotations.Pitch+dg.DPitch, 360),
		Yaw:   math.Mod(gyroRotations.Yaw+dg.DYaw, 360),
	}
}

func AccelerometerRotations(data types.XYZ) types.Rotations {
	roll := utils.RadToDeg(math.Atan2(data.Y, data.Z))
	pitch := -utils.RadToDeg(math.Atan2(data.X, data.Z))
	return types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	}
}

func applyFilter(pR float64, accR float64, gyroDR float64, lpfc float64) float64 {
	nR := (lpfc)*(pR+gyroDR) + (1-lpfc)*accR
	return math.Mod(nR, 360)
}

func CalcRotations(pR types.Rotations, aR types.Rotations, dg types.RotationsChanges, lowPassFilterCoefficient float64) types.Rotations {
	roll := applyFilter(pR.Roll, aR.Roll, dg.DRoll, lowPassFilterCoefficient)
	pitch := applyFilter(pR.Pitch, aR.Pitch, dg.DPitch, lowPassFilterCoefficient)
	// we don't use accelerometer yaw for the correction
	yaw := math.Mod(pR.Yaw+dg.DYaw, 360)
	return types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   yaw,
	}
}
