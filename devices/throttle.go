package devices

type throttleInput struct {
	input analogtodigital
	scale float64
}

func (t *throttleInput) Read() uint16 {
	digitalValue := t.input.Read()
	return uint16(float64(digitalValue) * t.scale)
}

func NewThrottle(input analogtodigital, digitalMaxValue uint16, maxValue uint16) *throttleInput {
	scale := float64(maxValue) / float64(digitalMaxValue)
	return &throttleInput{
		input: input,
		scale: scale,
	}
}
