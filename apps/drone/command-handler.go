package drone

import (
	"github.com/marksaravi/drone-go/apps/commons"
)

func (d *droneApp) applyCommands(commands []byte) {
	if d.offMotors(commands) {
		return
	}
	if d.onMotors(commands) {
		return
	}
	d.setCommands(commands)
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

func (d *droneApp) setCommands(commands []byte) {
	// roll := common.CalcRotationFromRawJoyStickRaw(commands[0:2], d.rollMidValue, d.rotationRange)
	// pitch := common.CalcRotationFromRawJoyStickRaw(commands[2:4], d.pitchlMidValue, d.rotationRange)
	// yaw := common.CalcRotationFromRawJoyStickRaw(commands[4:6], d.yawMidValue, d.rotationRange)
	throttle := commons.CalcThrottleFromRawJoyStickRaw(commands[6:8], d.maxThrottle)
	d.flightControl.flightState.SetThrottle(throttle)
	// if time.Since(d.lastImuPrint) >= time.Second/4 {
	// 	d.lastImuPrint = time.Now()
	// 	fmt.Println(throttle)
	// }

	// fmt.Printf("%6.2f, %6.2f, %6.2f, %6.2f \n", roll, pitch, yaw, throttle)
}
