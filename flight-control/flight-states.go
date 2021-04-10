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
	lastAcc  time.Time
	lastGyro time.Time
	millis   int
)

func init() {
	lastAcc = time.Now()
	lastGyro = time.Now()
	millis = 50
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

func lowPassFilter(curr types.Rotations, gyro types.XYZ, acc types.XYZ, dt float64, lowPassFilterCoefficient float64) types.Rotations {
	return types.Rotations{
		Roll:  utils.LowPassFilter(curr.Roll+gyro.X*dt, acc.X, lowPassFilterCoefficient),
		Pitch: utils.LowPassFilter(curr.Pitch+gyro.Y*dt, acc.Y, lowPassFilterCoefficient),
		Yaw:   utils.LowPassFilter(curr.Yaw+gyro.Z*dt, acc.Z, lowPassFilterCoefficient),
	}
}

func (fs *FlightStates) setRotations(lowPassFilterCoefficient float64) {
	fs.rotations = lowPassFilter(
		fs.rotations,
		fs.imuData.Gyro.Data,
		fs.imuData.Acc.Data,
		fs.imuData.Duration,
		lowPassFilterCoefficient,
	)
}

func (fs *FlightStates) ShowAccStates() {
	if time.Since(lastAcc) > time.Millisecond*time.Duration(millis) {
		ar := fs.accRotations
		fmt.Println(fmt.Sprintf("Acc: %.3f, %.3f, %.3f", ar.Roll, ar.Pitch, ar.Yaw))
		// fmt.Println(fmt.Sprintf("%.5f", ar.Roll))
		// fmt.Println(fmt.Sprintf("%.5f", ar.Pitch))
		// fmt.Println(fmt.Sprintf("%.5f", ar.Yaw))
		lastAcc = time.Now()
	}
}

func (fs *FlightStates) ShowGyroStates() {
	if time.Since(lastGyro) > time.Millisecond*time.Duration(millis) {
		gr := fs.gyroRotations
		fmt.Println(fmt.Sprintf("Gyr: %.3f, %.3f, %.3f", gr.Roll, gr.Pitch, gr.Yaw))
		// fmt.Println(fmt.Sprintf("%.5f", gr.Roll))
		// fmt.Println(fmt.Sprintf("%.5f", gr.Pitch))
		// fmt.Println(fmt.Sprintf("%.5f", gr.Yaw))
		lastGyro = time.Now()
	}
}
