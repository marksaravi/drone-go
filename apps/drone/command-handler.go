package drone

import (
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/utils"
)

const (
	COMMAND_TURN_ON     byte = 8
	COMMAND_TURN_OFF    byte = 128
	COMMAND_CALIB_INC_P byte = 1
	COMMAND_CALIB_DEC_P byte = 16
	COMMAND_CALIB_INC_I byte = 2
	COMMAND_CALIB_DEC_I byte = 32
	COMMAND_CALIB_INC_D byte = 4
	COMMAND_CALIB_DEC_D byte = 64
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
	// pressedButtons := commands[8]
	pushButtons := commands[9]

	d.onMotors(pushButtons)
	d.getThrottleCommands(hThrottle, lThrottle)
	d.getRotationCommands(hRoll, lRoll, hPitch, lPitch)

	// if d.flightControl.calibrationMode {
	// 	d.calibratePID(pressedButtons)
	// }
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
	// d.changePIDGain(pushButtonsmmand, COMMAND_CALIB_INC_P, COMMAND_CALIB_DEC_P, "P")
	// d.changePIDGain(pushButtonsmmand, COMMAND_CALIB_INC_I, COMMAND_CALIB_DEC_I, "I")
	// d.changePIDGain(pushButtonsmmand, COMMAND_CALIB_INC_D, COMMAND_CALIB_DEC_D, "D")
}

func (d *droneApp) changePIDGain(pushButtonsmmand, incCommand, decCommand byte, gain string) {
	// inc := float64(0)
	// if pushButtonsmmand == incCommand {
	// 	inc = 1
	// } else if pushButtonsmmand == decCommand {
	// 	inc = -1
	// } else {
	// 	return
	// }
	// switch gain {
	// case "P":
	// 	d.flightControl.arm_0_2_PID.UpdateGainP(inc * d.flightControl.calibrationIncP)
	// 	d.flightControl.arm_1_3_PID.UpdateGainP(inc * d.flightControl.calibrationIncP)
	// case "I":
	// 	d.flightControl.arm_0_2_PID.UpdateGainI(inc * d.flightControl.calibrationIncI)
	// 	d.flightControl.arm_1_3_PID.UpdateGainI(inc * d.flightControl.calibrationIncI)
	// case "D":
	// 	d.flightControl.arm_0_2_PID.UpdateGainD(inc * d.flightControl.calibrationIncD)
	// 	d.flightControl.arm_1_3_PID.UpdateGainD(inc * d.flightControl.calibrationIncD)
	// }
}
