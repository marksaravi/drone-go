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
	millis     int
	counter    int
	sampleRate int
)

func init() {
	lastPrint = time.Now()
	millis = 1000
	counter = 0
	sampleRate = 0
}

type FlightStates struct {
	imuData       imu.ImuData
	accRotations  types.Rotations
	gyroRotations types.Rotations
	rotations     types.Rotations
}

func (fs *FlightStates) Set(imuData imu.ImuData, config types.FlightConfig) {
	fs.imuData = imuData
	fs.setAccRotations(config.AccLowPassFilterCoefficient)
	fs.setGyroRotations()
	fs.setRotations(config.RotationsLowPassFilterCoefficient)
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

func (fs *FlightStates) setGyroRotations() {
	curr := fs.gyroRotations   // current rotations by gyro
	wg := fs.imuData.Gyro.Data // angular velocity
	dt := fs.imuData.Duration  // time interval
	fs.gyroRotations = types.Rotations{
		Roll:  curr.Roll - wg.X*dt,
		Pitch: curr.Pitch - wg.Y*dt,
		Yaw:   curr.Yaw - wg.Z*dt,
	}
}

func (fs *FlightStates) setRotations(lowPassFilterCoefficient float64) {
	curr := fs.rotations
	acc := fs.accRotations
	gyro := fs.imuData.Gyro.Data
	dt := fs.imuData.Duration
	roll := utils.LowPassFilter(curr.Roll+gyro.X*dt, acc.Roll, lowPassFilterCoefficient)
	pitch := utils.LowPassFilter(curr.Pitch+gyro.Y*dt, acc.Pitch, lowPassFilterCoefficient)

	fs.rotations = types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   fs.gyroRotations.Yaw,
	}
}

func toJson(r types.Rotations) string {
	return fmt.Sprintf("{\"roll\": %.3f, \"pitch\": %.3f, \"yaw\": %.3f}", r.Roll, r.Pitch, r.Yaw)
}

func (fs *FlightStates) ShowRotations(sensor string) string {
	var r types.Rotations
	var name string
	var json string

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
		sampleRate = counter
		counter = 0
		lastPrint = time.Now()
		fmt.Println(fmt.Sprintf("%s: %.3f, %.3f, %.3f, %d", name, r.Roll, r.Pitch, r.Yaw, sampleRate))
	}

	json = fmt.Sprintf("{\"acc\": %s, \"gyro\": %s, \"rotations\": %s, \"sampleRate\": %d}",
		toJson(fs.accRotations),
		toJson(fs.gyroRotations),
		toJson(fs.rotations),
		sampleRate,
	)
	return json
}
