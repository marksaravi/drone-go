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
	angle      float64
)

func init() {
	lastPrint = time.Now()
	lastUDP = time.Now()
	millis = 1000
	counter = 0
	sampleRate = 0
	angle = 0
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
	fs.setRotations()
}

func goDurToDt(d int64) float64 {
	return float64(d) / 1e9
}

func getOffset(offset float64, dt float64) float64 {
	return dt * offset
}

func accelerometerDataToRollPitch(data types.XYZ) (roll, pitch float64) {
	roll = utils.RadToDeg(math.Atan2(data.Y, data.Z))
	pitch = -utils.RadToDeg(math.Atan2(data.X, data.Z))
	return
}

func gyroscopeDataToRollPitchYawChange(wg types.XYZ, readingInterval int64) (
	float64, float64, float64) { // angular velocity
	dt := goDurToDt(readingInterval) // reading interval
	return wg.X * dt, wg.Y * dt, wg.X * dt
}

func (fs *FlightStates) setAccRotations(lowPassFilterCoefficient float64) {
	roll, pitch := accelerometerDataToRollPitch(fs.imuData.Acc.Data)
	fs.accRotations = types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	}
}

func (fs *FlightStates) setGyroRotations() {
	curr := fs.gyroRotations // current rotations by gyro
	_, dPitch, dYaw := gyroscopeDataToRollPitchYawChange(
		fs.imuData.Gyro.Data,
		fs.imuData.ReadingInterval,
	)

	angle += float64(math.Pi*2.0) / float64(fs.Config.ImuDataPerSecond)
	if angle > 2.0*math.Pi {
		angle = 0
	}
	fs.gyroRotations = types.Rotations{
		Roll:  math.Sin(angle) * 45,
		Pitch: curr.Pitch + dPitch,
		Yaw:   curr.Yaw + dYaw,
	}
}

func (fs *FlightStates) setRotations() {
	accCoeff := fs.Config.AccLowPassFilterCoefficient
	rotCoeff := fs.Config.RotationsLowPassFilterCoefficient
	curr := fs.rotations
	accNewRoll, accNewPitch := accelerometerDataToRollPitch(fs.imuData.Acc.Data)
	accRoll := utils.LowPassFilter(curr.Roll, accNewRoll, accCoeff)
	accPitch := utils.LowPassFilter(curr.Pitch, accNewPitch, accCoeff)
	gyroDRoll, gyroDPitch, gyroDYaw := gyroscopeDataToRollPitchYawChange(fs.imuData.Gyro.Data, fs.imuData.ReadingInterval)
	roll := utils.LowPassFilter(curr.Roll+gyroDRoll, accRoll, rotCoeff)
	pitch := utils.LowPassFilter(curr.Pitch+gyroDPitch, accPitch, rotCoeff)
	yaw := curr.Yaw + gyroDYaw
	fs.rotations = types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   yaw,
	}
}

func toJson(r types.Rotations) string {
	return fmt.Sprintf("[%.2f,%.2f,%.2f]", r.Roll, r.Pitch, r.Yaw)
}

func (fs *FlightStates) ImuDataToJson() string {
	return fmt.Sprintf("{\"a\":%s,\"g\":%s,\"r\":%s}",
		toJson(fs.accRotations),
		toJson(fs.gyroRotations),
		toJson(fs.rotations),
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
			fmt.Println(fmt.Sprintf("%s:%.3f,%.3f,%.3f,%d", name, r.Roll, r.Pitch, r.Yaw, sampleRate))
		}
		sampleRate = counter
		counter = 0
		lastPrint = time.Now()
	}
}

func (fs *FlightStates) Calibrate() {
	const CALIBRATION_TIME = 3
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
