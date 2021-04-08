package flightcontrol

import (
	"fmt"
	"math"
	"time"

	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
)

var (
	lastAcc    time.Time
	lastGyro   time.Time
	millis     int
	accScaler  float64
	gyroScaler float64
)

func init() {
	lastAcc = time.Now()
	lastGyro = time.Now()
	millis = 250
	gyroScaler = 0
	accScaler = 0
}

type FlightStates struct {
	imuData       imu.ImuData
	accRotations  types.Rotations
	gyroRotations types.Rotations
	rotations     types.Rotations
}

func (fs *FlightStates) Set(imuData imu.ImuData) {
	fs.imuData = imuData
	fs.setAccRotations(imuData)
	fs.setGyroRotations(imuData)
}

func (fs *FlightStates) setAccRotations(imuData imu.ImuData) {
	roll := math.Atan2(imuData.Acc.Data.X, imuData.Acc.Data.Z)
	pitch := math.Atan2(imuData.Acc.Data.Y, imuData.Acc.Data.Z)

	fs.accRotations = types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	}
}

func (fs *FlightStates) setRotations() {

}

func (fs *FlightStates) setGyroRotations(imuData imu.ImuData) {
	fs.gyroRotations.Roll = fs.gyroRotations.Roll + fs.imuData.Gyro.Data.X*fs.imuData.Duration
	fs.gyroRotations.Pitch = fs.gyroRotations.Pitch + fs.imuData.Gyro.Data.Y*fs.imuData.Duration
	fs.gyroRotations.Yaw = fs.gyroRotations.Yaw + fs.imuData.Gyro.Data.Z*fs.imuData.Duration
}

func (fs *FlightStates) ShowAccStates() {
	s := fs.accRotations.ToDeg().Scaler()
	if math.Abs(s-accScaler) > 1 && time.Since(lastAcc) > time.Millisecond*time.Duration(millis) {
		ar := fs.accRotations.ToDeg()
		fmt.Println(fmt.Sprintf("Acc: %.3f, %.3f, %.3f", ar.Roll, ar.Pitch, ar.Yaw))
		accScaler = s
		lastAcc = time.Now()
	}
}

func (fs *FlightStates) ShowGyroStates() {
	s := fs.gyroRotations.Scaler()
	if math.Abs(s-gyroScaler) > 1 && time.Since(lastGyro) > time.Millisecond*time.Duration(millis) {
		gr := fs.accRotations.ToDeg()
		fmt.Println(fmt.Sprintf("Gyr: %.3f, %.3f, %.3f", gr.Roll, gr.Pitch, gr.Yaw))
		gyroScaler = s
		lastGyro = time.Now()
	}
}
