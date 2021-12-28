package devices

type analogtodigital interface {
	Read() byte
}

type joystickInput struct {
	input analogtodigital
	value byte
}

func (js *joystickInput) Read() byte {
	js.value = js.input.Read()
	return js.value
}

func NewJoystick(ad analogtodigital) *joystickInput {
	return &joystickInput{
		input: ad,
	}
}
