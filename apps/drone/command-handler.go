package drone

func (d *droneApp) applyCommands(commands []byte) {
	d.offMotors(commands)
	d.onMotors(commands)
}

func (d *droneApp) onMotors(commands []byte) {
	if commands[5] == 1 {
		d.flightControl.flightState.ConnectThrottle()
	}
}

func (d *droneApp) offMotors(commands []byte) {
	if commands[5] == 16 {
		d.flightControl.flightState.DisconnectThrottle()
	}
}
