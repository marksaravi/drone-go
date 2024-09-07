package drone

import (
	"fmt"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/utils"
)

const (
	COMMAND_TURN_ON  byte = 1
	COMMAND_TURN_OFF byte = 16

	COMMAND_TURN_LEFT  byte = 1
	COMMAND_TURN_RIGHT byte = 2

	COMMAND_CALIB_INC_P byte = 2
	COMMAND_CALIB_INC_I byte = 4
	COMMAND_CALIB_INC_D byte = 8

	COMMAND_CALIB_DEC_P byte = 32
	COMMAND_CALIB_DEC_I byte = 64
	COMMAND_CALIB_DEC_D byte = 128
)

func (d *droneApp) applyCommands(commands []byte) {

	lRoll := commands[0]
	hRoll := commands[1]
	lPitch := commands[2]
	hPitch := commands[3]
	// lYaw := commands[4]
	// hYaw := commands[5]
	lThrottle := commands[6]
	hThrottle := commands[7]
	pressedButtons := commands[8]
	pushButtons := commands[9]

	d.onMotors(pushButtons)
	d.getThrottleCommands(hThrottle, lThrottle)
	d.getRotationCommands(hRoll, lRoll, hPitch, lPitch)
	d.setHeading(pressedButtons)

	d.calibratePID(pushButtons)
}

func (d *droneApp) setHeading(pressedButtons byte) {
	if pressedButtons&COMMAND_TURN_LEFT != 0 {
		d.flightControl.changeHeading(true)
	} else if pressedButtons&COMMAND_TURN_RIGHT != 0 {
		d.flightControl.changeHeading(false)
	}
}

func (d *droneApp) onMotors(pushButtons byte) {
	if pushButtons&COMMAND_TURN_ON > 0 && d.flightControl.getThrottle() < 3 {
		d.flightControl.turnOnMotors(true)
	} else if pushButtons&COMMAND_TURN_OFF > 0 {
		d.flightControl.turnOnMotors(false)
	}
}

func (d *droneApp) getThrottleCommands(hThrottle, lThrottle byte) {
	throttle := utils.CommandToThrottle(hThrottle, lThrottle, d.maxThrottle)
	d.flightControl.setThrottle(throttle)
}

func (d *droneApp) getRotationCommands(hRoll, lRoll, hPitch, lPitch byte) {
	roll := utils.CommandToRotation(hRoll, lRoll, d.rotationRange)
	pitch := utils.CommandToRotation(hPitch, lPitch, d.rotationRange)
	d.flightControl.setTargetRotations(imu.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	})
}

func (d *droneApp) calibratePID(pushButtonsmmand byte) {
	if !d.flightControl.arm_0_2_pid.IsCalibrationEnabled() &&
		!d.flightControl.arm_1_3_pid.IsCalibrationEnabled() &&
		!d.flightControl.yaw_pid.IsCalibrationEnabled() {
		return
	}

	t := rune('p')
	inc := true
	call := false
	if isP, incP := parsPidCalibrationCommand(pushButtonsmmand, COMMAND_CALIB_INC_P, COMMAND_CALIB_DEC_P); isP {
		t = 'p'
		inc = incP
		call = true
	} else if isI, incI := parsPidCalibrationCommand(pushButtonsmmand, COMMAND_CALIB_INC_I, COMMAND_CALIB_DEC_I); isI {
		t = 'i'
		inc = incI
		call = true
	} else if isD, incD := parsPidCalibrationCommand(pushButtonsmmand, COMMAND_CALIB_INC_D, COMMAND_CALIB_DEC_D); isD {
		t = 'd'
		inc = incD
		call = true
	}
	if call {
		if d.flightControl.arm_0_2_pid.IsCalibrationEnabled() {
			d.flightControl.arm_0_2_pid.Calibrate(t, inc)
		}
		if d.flightControl.arm_1_3_pid.IsCalibrationEnabled() {
			d.flightControl.arm_1_3_pid.Calibrate(t, inc)
		}
		if d.flightControl.yaw_pid.IsCalibrationEnabled() {
			d.flightControl.yaw_pid.Calibrate(t, inc)
		}
	}
}

func parsPidCalibrationCommand(pushButtonsmmand, incCommand, decCommand byte) (is, inc bool) {
	is = pushButtonsmmand&incCommand != 0 || pushButtonsmmand&decCommand != 0
	inc = pushButtonsmmand&incCommand != 0
	if is {
		fmt.Println(is, inc)
	}
	return
}
