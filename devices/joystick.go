package devices

type analogtodigital interface {
	Read() float32
}

type joystickInput struct {
	input analogtodigital
	value float32
}

func (js *joystickInput) Read() float32 {
	js.value = js.input.Read()
	return js.value
}

func NewJoystick(ad analogtodigital) *joystickInput {
	return &joystickInput{
		input: ad,
	}
}
