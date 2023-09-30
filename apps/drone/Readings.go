package drone

import (
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

func (d *droneApp) ReadIMU() (imu.Rotations, bool) {
	if time.Since(d.lastImuData) < time.Second/time.Duration(d.imuDataPerSecond) {
		return imu.Rotations{}, false
	}
	d.lastImuData = time.Now()
	rotations, err := d.imu.Read()
	if err != nil {
		return rotations, false
	}
	return rotations, true
}

func (d *droneApp) ReceiveCommand() ([]byte, bool) {
	if time.Since(d.lastCommand) < time.Second/time.Duration(2*d.commandsPerSecond) {
		return nil, false
	}
	d.lastCommand = time.Now()
	return d.receiver.Receive()
}
