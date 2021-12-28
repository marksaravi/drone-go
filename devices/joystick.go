package devices

type analogtodigital interface {
	Read() uint16
}

type joystickInput struct {
	input analogtodigital
	value uint16
}

func (js *joystickInput) Read() uint16 {
	js.value = js.input.Read()
	return js.value
}

func NewJoystick(input analogtodigital) *joystickInput {
	return &joystickInput{
		input: input,
	}
}
