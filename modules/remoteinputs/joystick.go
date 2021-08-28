package remoteinputs

type joystick interface {
	Read() float32
}

type joystickInput struct {
	input     joystick
	value     float32
	isChanged bool
}

func (js *joystickInput) read() {
	pv := js.value
	js.value = js.input.Read()
	js.isChanged = js.value != pv
}

func (ri *remoteInputs) readJoysticks() {
	ri.roll.read()
	ri.pitch.read()
	ri.yaw.read()
}
