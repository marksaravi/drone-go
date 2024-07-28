package drone

import (
	"fmt"

	"github.com/marksaravi/drone-go/apps/commons"
	"github.com/marksaravi/drone-go/devices/imu"
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
	d.onMotors(commands)
	d.getThrottleCommands(commands)
	d.getRotationCommands(commands)
	if d.flightControl.calibrationMode {
		d.calibratePID(commands[9])
	}
}

func (d *droneApp) onMotors(commands []byte) {
	if commands[9] == COMMAND_TURN_ON {
		d.flightControl.turnOnMotors(true)
	} else if commands[9] == COMMAND_TURN_OFF {
		d.flightControl.turnOnMotors(false)
	}
}

func (d *droneApp) getThrottleCommands(commands []byte) {
	throttle := commons.CalcThrottleFromRawJoyStickRaw(commands[6:8], d.maxThrottle, d.throttleZeroOffset)
	d.flightControl.setThrottle(throttle)
}

func (d *droneApp) getRotationCommands(commands []byte) {
	roll := commons.CalcRotationFromRawJoyStickRaw(commands[0:2], d.rollMidValue, d.rotationRange)
	pitch := commons.CalcRotationFromRawJoyStickRaw(commands[2:4], d.pitchlMidValue, d.rotationRange)
	// yaw := commons.CalcRotationFromRawJoyStickRaw(commands[4:6], d.yawMidValue, d.rotationRange)
	d.flightControl.setTargetRotations(imu.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	})
}
func (d *droneApp) calibratePID(command byte) {
	if command == 0 {
		return
	}
	d.changePIDGain(command, COMMAND_CALIB_INC_P, COMMAND_CALIB_DEC_P, "P")
	d.changePIDGain(command, COMMAND_CALIB_INC_I, COMMAND_CALIB_DEC_I, "I")
	d.changePIDGain(command, COMMAND_CALIB_INC_D, COMMAND_CALIB_DEC_D, "D")
}

func (d *droneApp) changePIDGain(command, incCommand, decCommand byte, gain string) {
	inc := float64(0)
	if command == incCommand {
		inc = 1
	} else if command == decCommand {
		inc = -1
	}
	if inc == 0 {
		return
	}
	switch gain {
	case "P":
		d.flightControl.arm_0_2_PID.UpdateGainP(inc * d.flightControl.calibrationIncP)
		d.flightControl.arm_1_3_PID.UpdateGainP(inc * d.flightControl.calibrationIncP)
		fmt.Println("P: ", d.flightControl.arm_0_2_PID.GainP(), d.flightControl.arm_1_3_PID.GainP())
	case "I":
		d.flightControl.arm_0_2_PID.UpdateGainI(inc * d.flightControl.calibrationIncI)
		d.flightControl.arm_1_3_PID.UpdateGainI(inc * d.flightControl.calibrationIncI)
		fmt.Println("I: ", d.flightControl.arm_0_2_PID.GainI(), d.flightControl.arm_1_3_PID.GainI())
	case "D":
		d.flightControl.arm_0_2_PID.UpdateGainD(inc * d.flightControl.calibrationIncD)
		d.flightControl.arm_1_3_PID.UpdateGainD(inc * d.flightControl.calibrationIncD)
		fmt.Println("D: ", d.flightControl.arm_0_2_PID.GainD(), d.flightControl.arm_1_3_PID.GainD())
	}
}
