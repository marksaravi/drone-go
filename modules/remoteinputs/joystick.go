package remoteinputs

type joystick interface {
	Read() float32
}

type joystickInput struct {
	input joystick
	value float32
}

func (js *joystickInput) read() (isChanged bool) {
	isChanged = false
	pv := js.value
	js.value = js.input.Read()
	isChanged = js.value != pv
	return
}

func (ri *remoteInputs) readJoysticks() (isChanged bool) {
	rollChanged := ri.roll.read()
	pitchChanged := ri.pitch.read()
	yawChanged := ri.yaw.read()
	return rollChanged || pitchChanged || yawChanged
}
