package devices

type analogtodigital interface {
	Read() uint16
}

type joystickInput struct {
	input    analogtodigital
	offset   int
	dir      int
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
	return digitalValue * js.dir
}

func NewJoystick(input analogtodigital, maxValue int, offset int, dir int) *joystickInput {
	return &joystickInput{
		input:    input,
		maxValue: maxValue,
		offset:   offset,
		dir:      dir,
	}
}
