package remotecontrol

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type button interface {
	Read() models.ButtonData
}

type joystick interface {
	Read() models.JoystickData
}

type remoteControl struct {
	roll         joystick
	pitch        joystick
	yaw          joystick
	btnFrontLeft button
	data         models.RemoteControlData
}

func NewRemoteControl(roll joystick, pitch joystick, yaw joystick, btnFrontLeft button) *remoteControl {
	return &remoteControl{
		roll:         roll,
		pitch:        pitch,
		yaw:          yaw,
		btnFrontLeft: btnFrontLeft,
	}
}

func (rc *remoteControl) Start() {
	for {
		rc.read()
		rc.showData()
		time.Sleep(250 * time.Millisecond)
	}
}

func (rc *remoteControl) read() {
	rc.data = models.RemoteControlData{
		Roll:            rc.roll.Read(),
		Pitch:           rc.pitch.Read(),
		Yaw:             rc.yaw.Read(),
		ButtonFrontLeft: rc.btnFrontLeft.Read(),
	}
}

func (rc *remoteControl) showData() {
	fmt.Println(rc.data)
}
