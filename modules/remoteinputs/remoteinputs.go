package remoteinputs

import "fmt"

type remoteInputs struct {
	roll            *joystickInput
	pitch           *joystickInput
	yaw             *joystickInput
	buttonFrontLeft *buttonInput
	stopped         bool
}

func NewRemoteInputs(roll joystick, pitch joystick, yaw joystick, inputButtonFrontLeft button) *remoteInputs {
	return &remoteInputs{
		roll: &joystickInput{
			input: roll,
			value: 0,
		},
		pitch: &joystickInput{
			input: pitch,
			value: 0,
		},
		yaw: &joystickInput{
			input: yaw,
			value: 0,
		},
		buttonFrontLeft: &buttonInput{
			input: inputButtonFrontLeft,
			value: false,
		},
		stopped: false,
	}
}

func (ri *remoteInputs) ReadInputs() (isStopChanged bool) {
	isStopChanged = ri.readStopButtons()
	ri.readJoysticks()
	return
}

func (ri *remoteInputs) IsStopped() bool {
	return ri.stopped
}

func (ri *remoteInputs) PrintData() {
	fmt.Printf("Button-Front-Left: %t, ", ri.buttonFrontLeft.value)
	fmt.Printf("roll: %5.2f, ", ri.roll.value)
	fmt.Printf("pitch: %5.2f, ", ri.pitch.value)
	fmt.Printf("yaw: %5.2f", ri.yaw.value)
	fmt.Println()
}
