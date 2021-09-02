package remotecontrol

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type radio interface {
	IsDataAvailable() bool
	ReceiverOn()
	ReceiveFlightData() models.FlightData
	TransmitterOn()
	TransmitFlightData(models.FlightData) error
}

type button interface {
	Read() models.ButtonData
}

type joystick interface {
	Read() models.JoystickData
}

type remoteControl struct {
	radio        radio
	roll         joystick
	pitch        joystick
	yaw          joystick
	throttle     joystick
	btnFrontLeft button
	data         models.RemoteControlData
}

func NewRemoteControl(radio radio, roll, pitch, yaw, throttle joystick, btnFrontLeft button) *remoteControl {
	return &remoteControl{
		radio:        radio,
		roll:         roll,
		pitch:        pitch,
		yaw:          yaw,
		throttle:     throttle,
		btnFrontLeft: btnFrontLeft,
	}
}

func (rc *remoteControl) Start() {
	sendTimer := time.Tick(time.Second / 25)
	rc.radio.ReceiverOn()
	var id uint32 = 0
	for range sendTimer {
		rc.read()
		rc.showData()
		rc.radio.TransmitterOn()
		rc.radio.TransmitFlightData(models.FlightData{
			Id:              id,
			Roll:            rc.data.Roll.Value,
			Pitch:           rc.data.Pitch.Value,
			Yaw:             rc.data.Yaw.Value,
			Throttle:        rc.data.Throttle.Value,
			Altitude:        0,
			IsRemoteControl: true,
			IsDrone:         false,
			IsMotorsEngaged: false,
		})
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

var lastPrint time.Time = time.Now()

func (rc *remoteControl) showData() {
	if time.Since(lastPrint) < time.Second/4 {
		return
	}
	lastPrint = time.Now()
	fmt.Println(rc.data)
}
