package remotecontrol

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type radio interface {
	ReceiverOn()
	ReceiveFlightData() (models.FlightData, bool)
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

func (rc *remoteControl) Start() {
	rc.radio.ReceiverOn()
	acknowleg := createAckReceiver(rc.radio)

	sendTimer := time.NewTicker(time.Second / 25)
	var id uint32 = 0
	lastAcknowleged := time.Now()
	var flightData models.FlightData = models.FlightData{
		Id: 0,
	}
	for {
		select {
		case <-sendTimer.C:
			rc.read()
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
			rc.radio.ReceiverOn()
			id++
		case flightData = <-acknowleg:
			lastAcknowleged = time.Now()
		default:
			if time.Since(lastAcknowleged) > time.Millisecond*200 {
				fmt.Println("Connection Error ", flightData.Id)
			}
		}
	}
}

func createAckReceiver(receiver radio) <-chan models.FlightData {
	acknowlegChannel := make(chan models.FlightData)
	go func(receiver radio, ackChannell chan models.FlightData) {
		for {
			ack, isavailable := receiver.ReceiveFlightData()
			if isavailable {
				ackChannell <- ack
			}

		}
	}(receiver, acknowlegChannel)
	return acknowlegChannel
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

func (rc *remoteControl) showData(id uint32) {
	if time.Since(lastPrint) < time.Second/4 {
		return
	}
	lastPrint = time.Now()
	fmt.Println(id, rc.data)
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
