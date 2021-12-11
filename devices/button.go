package devices

type gpioswitch interface {
	Read() bool
}

type button struct {
	input gpioswitch
	value bool
}

func (btn *button) Read() bool {
	btn.value = btn.input.Read()
	return btn.value
}

func NewButton(input gpioswitch) *button {
	return &button{
		input: input,
	}
}
