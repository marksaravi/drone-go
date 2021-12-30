package devices

type throttleInput struct {
	input analogtodigital
	value int
}

func (js *throttleInput) Read() int {
	digitalValue := js.input.Read()
	return int(digitalValue)
}

func NewThrottleInput(input analogtodigital) *throttleInput {
	return &throttleInput{
		input: input,
	}
}
