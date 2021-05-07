package flightcontrol

import (
	"fmt"
	"math"
	"time"

	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

var (
	lastPrint  time.Time
	lastUDP    time.Time
	millis     int
	counter    int
	sampleRate int
)

func init() {
	lastPrint = time.Now()
	lastUDP = time.Now()
	millis = 1000
	counter = 0
	sampleRate = 0
}

type FlightStates struct {
	Config         types.FlightConfig
	ImuDataChannel <-chan imu.ImuData
	imuData        imu.ImuData
	accRotations   types.Rotations
	gyroRotations  types.Rotations
	rotations      types.Rotations
}

func (fs *FlightStates) Reset() {
	fs.gyroRotations = types.Rotations{
		Roll:  0,
		Pitch: 0,
		Yaw:   0,
	}
}

func (fs *FlightStates) Set(imuData imu.ImuData) {
	fs.imuData = imuData
	fs.setAccRotations(fs.Config.AccLowPassFilterCoefficient)
	fs.setGyroRotations()
	fs.setRotations(fs.Config.RotationsLowPassFilterCoefficient)
}

func (fs *FlightStates) setAccRotations(lowPassFilterCoefficient float64) {
	x := fs.imuData.Acc.Data.X
	y := fs.imuData.Acc.Data.Y
	z := fs.imuData.Acc.Data.Z
	roll := utils.RadToDeg(math.Atan2(y, z))
	pitch := utils.RadToDeg(math.Atan2(-x, math.Sqrt(y*y+z*z)))

	fs.accRotations = types.Rotations{
		Roll:  utils.LowPassFilter(fs.accRotations.Roll, roll, lowPassFilterCoefficient),
		Pitch: utils.LowPassFilter(fs.accRotations.Pitch, pitch, lowPassFilterCoefficient),
		Yaw:   0,
	}
}

func goDurToDt(d int64) float64 {
	return float64(d) / 1e9
}

func getOffset(offset float64, dt float64) float64 {
	return dt * offset
}

func (fs *FlightStates) setGyroRotations() {
	curr := fs.gyroRotations                    // current rotations by gyro
	wg := fs.imuData.Gyro.Data                  // angular velocity
	dt := goDurToDt(fs.imuData.ReadingInterval) // reading interval
	fs.gyroRotations = types.Rotations{
		Roll:  curr.Roll - wg.X*dt - getOffset(fs.Config.GyroscopeOffsets.X, dt),
		Pitch: curr.Pitch - wg.Y*dt - getOffset(fs.Config.GyroscopeOffsets.Y, dt),
		Yaw:   curr.Yaw - wg.Z*dt - getOffset(fs.Config.GyroscopeOffsets.Z, dt),
	}
}

func (fs *FlightStates) setRotations(lowPassFilterCoefficient float64) {
	curr := fs.rotations
	acc := fs.accRotations
	wg := fs.imuData.Gyro.Data // angular velocity
	dt := goDurToDt(fs.imuData.ReadingInterval)
	roll := utils.LowPassFilter(curr.Roll-wg.X*dt, acc.Roll, lowPassFilterCoefficient)
	pitch := utils.LowPassFilter(curr.Pitch-wg.Y*dt, acc.Pitch, lowPassFilterCoefficient)

	fs.rotations = types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   fs.gyroRotations.Yaw,
	}
}

func toJson(r types.Rotations) string {
	return fmt.Sprintf("{\"roll\": %.3f, \"pitch\": %.3f, \"yaw\": %.3f}", r.Roll, r.Pitch, r.Yaw)
}

func (fs *FlightStates) ImuDataToJson() string {
	return fmt.Sprintf("{\"accelerometer\": %s, \"gyroscope\": %s, \"rotations\": %s, \"readingInterval\": %d, \"elapsedTime\": %d, \"sampleRate\": %d, \"totalData\": %d}",
		toJson(fs.accRotations),
		toJson(fs.gyroRotations),
		toJson(fs.rotations),
		fs.imuData.ReadingInterval,
		fs.imuData.ElapsedTime,
		fs.imuData.SampleRate,
		fs.imuData.TotalData,
	)
}

func (fs *FlightStates) ShowRotations(sensor string, json string) {
	var r types.Rotations
	var name string

	switch sensor {
	case "acc":
		r = fs.accRotations
		name = "Acc"
	case "gyro":
		r = fs.gyroRotations
		name = "Gyro"
	default:
		r = fs.rotations
		name = "Rotations"
	}
	counter++

	if time.Since(lastPrint) > time.Millisecond*time.Duration(millis) {
		if sensor == "json" {
			fmt.Println(json)
		} else {
			fmt.Println(fmt.Sprintf("%s: %.3f, %.3f, %.3f, %d", name, r.Roll, r.Pitch, r.Yaw, sampleRate))
		}
		sampleRate = counter
		counter = 0
		lastPrint = time.Now()
	}
}

func (fs *FlightStates) Calibrate() {
	const CALIBRATION_TIME = 10
	fs.Reset()
	fmt.Println("Calibration started...")
	start := time.Now()
	for time.Since(start) < time.Second*CALIBRATION_TIME {
		imuData := <-fs.ImuDataChannel
		fs.Set(imuData)
	}
	fs.Config.GyroscopeOffsets = types.Offsets{
		X: fs.gyroRotations.Roll / CALIBRATION_TIME,
		Y: fs.gyroRotations.Pitch / CALIBRATION_TIME,
		Z: fs.gyroRotations.Yaw / CALIBRATION_TIME,
	}
	fmt.Println("Calibration finished.", fs.Config.GyroscopeOffsets)
}
