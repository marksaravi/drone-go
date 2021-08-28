package remoteinputs

import "fmt"

type button interface {
	Read() bool
}

type joystick interface {
	Read() float32
}

type joystickInput struct {
	input joystick
	value float32
}

type remoteInputs struct {
	roll                 *joystickInput
	inputButtonFrontLeft button

	valueButtonFrontLeft bool
	stopped              bool
}

func NewRemoteInputs(roll joystick, inputButtonFrontLeft button) *remoteInputs {
	return &remoteInputs{
		roll: &joystickInput{
			input: roll,
			value: 0,
		},
		inputButtonFrontLeft: inputButtonFrontLeft,
		stopped:              false,
	}
}

func (ri *remoteInputs) RefreshInputs() (isStopChanged bool) {
	isStopChanged = ri.readStopButtons()
	ri.readJoysticks()
	return
}

func readJoystick(js *joystickInput) (isChanged bool) {
	isChanged = false
	pv := js.value
	js.value = js.input.Read()
	isChanged = js.value != pv
	return
}

func (ri *remoteInputs) readJoysticks() (isChanged bool) {
	rollChanged := readJoystick(ri.roll)
	return rollChanged
}

func (ri *remoteInputs) readStopButtons() (isChanged bool) {
	isChanged = false
	ri.valueButtonFrontLeft = ri.inputButtonFrontLeft.Read()
	if ri.valueButtonFrontLeft {
		if !ri.stopped {
			isChanged = true
		}
		ri.stopped = true
	}
	return isChanged
}

func (ri *remoteInputs) IsStopped() bool {
	return ri.stopped
}

func (ri *remoteInputs) PrintData() {
	fmt.Printf("Button-Front-Left: %t, ", ri.valueButtonFrontLeft)
	fmt.Printf("roll: %f", ri.roll.value)
	fmt.Println()
}
