package remoteinputs

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/models"
)

type remoteInputs struct {
	roll    *joystickInput
	pitch   *joystickInput
	yaw     *joystickInput
	stopped bool
}

func NewRemoteInputs(roll joystick, pitch joystick, yaw joystick) *remoteInputs {
	return &remoteInputs{
		roll: &joystickInput{
			input: roll,
		},
		pitch: &joystickInput{
			input: pitch,
		},
		yaw: &joystickInput{
			input: yaw,
		},
		stopped: false,
	}
}

func (ri *remoteInputs) ReadInputs() models.RemoteControlData {
	ri.readJoysticks()
	return models.RemoteControlData{
		Roll:  ri.roll.data,
		Pitch: ri.pitch.data,
		Yaw:   ri.yaw.data,
	}
}

func (ri *remoteInputs) PrintData() {
	fmt.Printf("roll: %5.2f, ", ri.roll.data.Value)
	fmt.Printf("pitch: %5.2f, ", ri.pitch.data.Value)
	fmt.Printf("yaw: %5.2f", ri.yaw.data.Value)
	fmt.Println()
}
