package devices

type analogtodigital interface {
	Read() uint16
}

type joystickInput struct {
	input    analogtodigital
	offset   int
	maxValue int
}

func (js *joystickInput) Read() int {
	digitalValue := int(js.input.Read()) - js.offset
	if digitalValue > (js.maxValue - 1) {
		return js.maxValue - 1
	}
	if digitalValue < 0 {
		return 0
	}
	return digitalValue
}

func NewJoystick(input analogtodigital, maxValue int, offset int) *joystickInput {
	return &joystickInput{
		input:    input,
		maxValue: maxValue,
		offset:   offset,
	}
}
