package remotecontrol

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radio interface {
	ReceiverOn()
	Receive() ([32]byte, bool)
	TransmitterOn()
	Transmit([32]byte) error
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
	data         models.FlightCommands
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
			fc := models.FlightCommands{
				Id:                id,
				Roll:              rc.data.Roll,
				Pitch:             rc.data.Pitch,
				Yaw:               rc.data.Yaw,
				Throttle:          rc.data.Throttle,
				ButtonFrontLeft:   rc.data.ButtonFrontLeft,
				ButtonFrontRight:  rc.data.ButtonFrontRight,
				ButtonTopLeft:     rc.data.ButtonTopLeft,
				ButtonTopRight:    rc.data.ButtonTopRight,
				ButtonBottomLeft:  rc.data.ButtonBottomLeft,
				ButtonBottomRight: rc.data.ButtonBottomRight,
			}
			fmt.Printf(
				"%0.2f, %0.2f, %0.2f, %0.2f, %v, %v, %v, %v, %v, %v\n",
				fc.Roll,
				fc.Pitch,
				fc.Yaw, fc.Throttle,
				fc.ButtonFrontLeft,
				fc.ButtonFrontRight,
				fc.ButtonTopLeft,
				fc.ButtonTopRight, fc.ButtonBottomLeft,
				fc.ButtonBottomRight,
			)
			payload := utils.SerializeFlightCommand(fc)
			rc.radio.Transmit(
				payload)
			rc.radio.ReceiverOn()
			id++
		case flightCommands = <-acknowleg:
			lastAcknowleged = time.Now()
		default:
			if time.Since(lastAcknowleged) > time.Millisecond*1000000 {
				fmt.Println("Connection Error ", flightCommands.Id)
			}
		}
	}
}

func createAckReceiver(receiver radio) <-chan models.FlightCommands {
	acknowlegChannel := make(chan models.FlightCommands)
	go func(receiver radio, ackChannell chan models.FlightCommands) {
		for {
			ack, isavailable := receiver.Receive()
			if isavailable {
				ackChannell <- utils.DeserializeFlightCommand(ack)
			}

		}
	}(receiver, acknowlegChannel)
	return acknowlegChannel
}

func (rc *remoteControl) read() {
	rc.data = models.FlightCommands{
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
