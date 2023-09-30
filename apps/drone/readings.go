package drone

import (
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

func (d *droneApp) ReadIMU() bool {
	if time.Since(d.lastImuData) < time.Second/time.Duration(d.imuDataPerSecond) {
		return false
	}
	d.lastImuData = time.Now()
	rotations, accRotations, gyroRotations, err := d.imu.Read()
	if err != nil {
		return false
	}
	d.rotations = rotations
	d.accRotations = accRotations
	d.gyroRotations = gyroRotations
	return true
}

func (d *droneApp) ReceiveCommand() ([]byte, bool) {
	if time.Since(d.lastCommand) < time.Second/time.Duration(2*d.commandsPerSecond) {
		return nil, false
	}
	d.lastCommand = time.Now()
	return d.receiver.Receive()
}

func (d *droneApp) PlotterData() {
	if !d.plotterActive {
		return
	}
}

func SerializeRotations(r imu.Rotations) {
}
