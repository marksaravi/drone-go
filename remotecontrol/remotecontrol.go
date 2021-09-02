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
	throttle     joystick
	btnFrontLeft button
	data         models.RemoteControlData
}

func NewRemoteControl(roll, pitch, yaw, throttle joystick, btnFrontLeft button) *remoteControl {
	return &remoteControl{
		roll:         roll,
		pitch:        pitch,
		yaw:          yaw,
		throttle:     throttle,
		btnFrontLeft: btnFrontLeft,
	}
}

func (rc *remoteControl) Start() {
	sendTimer := time.Tick(time.Second / 25)
	for range sendTimer {
		rc.read()
		rc.showData()
	}
}

func (rc *remoteControl) read() {
	rc.data = models.RemoteControlData{
		Roll:            rc.roll.Read(),
		Pitch:           rc.pitch.Read(),
		Yaw:             rc.yaw.Read(),
		Throttle:        rc.throttle.Read(),
		ButtonFrontLeft: rc.btnFrontLeft.Read(),
	}
}

func (rc *remoteControl) showData() {
	fmt.Println(rc.data)
}
