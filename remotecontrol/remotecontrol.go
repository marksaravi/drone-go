package remotecontrol

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type radio interface {
	ReceiverOn()
	ReceiveFlightData() (models.FlightCommands, bool)
	TransmitterOn()
	TransmitFlightData(models.FlightCommands) error
}

type button interface {
	Read() bool
}

type joystick interface {
	Read() float32
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
	var flightCommands models.FlightCommands = models.FlightCommands{
		Id: 0,
	}
	for {
		select {
		case <-sendTimer.C:
			rc.read()
			rc.radio.TransmitterOn()
			rc.radio.TransmitFlightData(models.FlightCommands{
				Id:              id,
				Roll:            rc.data.Roll,
				Pitch:           rc.data.Pitch,
				Yaw:             rc.data.Yaw,
				Throttle:        rc.data.Throttle,
				Altitude:        0,
				IsRemoteControl: true,
				IsDrone:         false,
				IsMotorsEngaged: false,
			})
			rc.radio.ReceiverOn()
			id++
		case flightCommands = <-acknowleg:
			lastAcknowleged = time.Now()
		default:
			if time.Since(lastAcknowleged) > time.Millisecond*200 {
				fmt.Println("Connection Error ", flightCommands.Id)
			}
		}
	}
}

func createAckReceiver(receiver radio) <-chan models.FlightCommands {
	acknowlegChannel := make(chan models.FlightCommands)
	go func(receiver radio, ackChannell chan models.FlightCommands) {
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
