package drone

import (
	"fmt"

	"github.com/marksaravi/drone-go/apps/commons"
	"github.com/marksaravi/drone-go/devices/imu"
)

func (d *droneApp) applyCommands(commands []byte) {
	d.onMotors(commands)
	d.getThrottleCommands(commands)
	d.getRotationCommands(commands)
	if d.flightControl.calibrationMode {
		d.calibratePID(commands[9])
	}
}

func (d *droneApp) onMotors(commands []byte) {
	if commands[9] == 1 {
		d.flightControl.turnOnMotors(true)
	} else if commands[9] == 16 {
		d.flightControl.turnOnMotors(false)
	}
}

func (d *droneApp) getThrottleCommands(commands []byte) {
	throttle := commons.CalcThrottleFromRawJoyStickRaw(commands[6:8], d.maxThrottle)
	d.flightControl.SetThrottle(throttle)
}

func (d *droneApp) getRotationCommands(commands []byte) {
	roll := commons.CalcRotationFromRawJoyStickRaw(commands[0:2], d.rollMidValue, d.rotationRange)
	pitch := commons.CalcRotationFromRawJoyStickRaw(commands[2:4], d.pitchlMidValue, d.rotationRange)
	// yaw := commons.CalcRotationFromRawJoyStickRaw(commands[4:6], d.yawMidValue, d.rotationRange)
	d.flightControl.SetTargetRotations(imu.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	})
}
func (d *droneApp) calibratePID(command byte) {
	if command == 0 {
		return
	}
	pInc := d.changePIDGain(command, 2, 32)
	iInc := d.changePIDGain(command, 4, 64)
	dInc := d.changePIDGain(command, 8, 128)
	fmt.Println(pInc, iInc, dInc, command)
}

func (d *droneApp) changePIDGain(command, incCommand, decCommand byte) int {
	if command == incCommand {
		return 1
	} else if command == decCommand {
		return -1
	}
	return 0
}
