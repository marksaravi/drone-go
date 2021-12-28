package devices

type analogtodigital interface {
	Read() int
}

type joystickInput struct {
	input analogtodigital
	value int
}

func (js *joystickInput) Read() int {
	js.value = js.input.Read()
	return js.value
}

func NewJoystick(input analogtodigital) *joystickInput {
	return &joystickInput{
		input: input,
	}
}
