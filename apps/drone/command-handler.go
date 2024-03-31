package drone

import (
	"github.com/marksaravi/drone-go/apps/commons"
	"github.com/marksaravi/drone-go/devices/imu"
)

func (d *droneApp) applyCommands(commands []byte) {
	if d.offMotors(commands) {
		return
	}
	if d.onMotors(commands) {
		return
	}
	d.getThrottleCommands(commands)
	d.getRotationCommands(commands)
}

func (d *droneApp) onMotors(commands []byte) bool {
	if commands[9] == 1 {
		d.flightControl.SetToZeroThrottleState(true)
		return true
	}
	return false
}

func (d *droneApp) offMotors(commands []byte) bool {
	if commands[9] == 16 {
		d.flightControl.SetToZeroThrottleState(false)
		return true
	}
	return false
}

func (d *droneApp) getThrottleCommands(commands []byte) {
	throttle := commons.CalcThrottleFromRawJoyStickRaw(commands[6:8], d.maxThrottle)
	d.flightControl.flightState.SetThrottle(throttle)
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
