package flightcontrol

import (
	"fmt"
	"math"
	"time"

	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/utils/euler"
)

var (
	lastPrint           time.Time
	gyroX, gyroY, gyroZ float64
	currValue           float64
)

func init() {
	lastPrint = time.Now()
	gyroX = 0
	gyroY = 0
	gyroZ = 0
	currValue = 1000
}

type FlightStates struct {
	imuData imu.ImuData
}

func (fs *FlightStates) SetImuData(imuData imu.ImuData) {
	fs.imuData = imuData
}

func (fs *FlightStates) ShowStates() {
	gyroX = gyroX + fs.imuData.Gyro.Data.X*fs.imuData.Duration
	gyroY = gyroY + fs.imuData.Gyro.Data.Y*fs.imuData.Duration
	gyroZ = gyroZ + fs.imuData.Gyro.Data.Z*fs.imuData.Duration
	x := fs.imuData.Acc.Data.X
	y := fs.imuData.Acc.Data.Y
	z := fs.imuData.Acc.Data.Z

	v := math.Sqrt(x*x + y*y + z*z)

	if math.Abs(currValue-v) > 0.025 && time.Since(lastPrint) > time.Millisecond*250 {
		e, _ := euler.AccelerometerToEulerAngles(fs.imuData.Acc.Data)
		fmt.Println(fmt.Sprintf("%.3f, %.3f, %.3f, %.3f, %.3f", x, y, z, e.Theta, e.Phi))
		lastPrint = time.Now()
		currValue = v
	}
}
