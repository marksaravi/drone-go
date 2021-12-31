package remotecontrol

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/devices/radio"
	piezobuzzer "github.com/marksaravi/drone-go/hardware/piezo-buzzer"
	"github.com/marksaravi/drone-go/models"
)

type button interface {
	Read() bool
}

type joystick interface {
	Read() int
}

type remoteControl struct {
	commandPerSecond                int
	radio                           models.Radio
	roll                            joystick
	pitch                           joystick
	yaw                             joystick
	throttle                        joystick
	btnFrontLeft                    button
	btnFrontRight                   button
	btnTopLeft                      button
	btnTopRight                     button
	btnBottomLeft                   button
	btnBottomRight                  button
	buzzer                          *piezobuzzer.Buzzer
	connectionState                 radio.ConnectionState
	shutdownCountdown               time.Time
	suppressLostConnectionCountdown time.Time
}

func (rc *remoteControl) read() models.FlightCommands {
	roll := rc.roll.Read()
	pitch := rc.pitch.Read()
	yaw := rc.yaw.Read()
	throttle := rc.throttle.Read()
	buttonFrontLeft := rc.btnFrontLeft.Read()
	buttonFrontRight := rc.btnFrontRight.Read()
	buttonTopLeft := rc.btnTopLeft.Read()
	buttonTopRight := rc.btnTopRight.Read()
	buttonBottomLeft := rc.btnBottomLeft.Read()
	buttonBottomRight := rc.btnBottomRight.Read()

	return models.FlightCommands{
		Roll:              uint16(roll),
		Pitch:             uint16(pitch),
		Yaw:               uint16(yaw),
		Throttle:          uint16(throttle),
		ButtonFrontLeft:   buttonFrontLeft,
		ButtonFrontRight:  buttonFrontRight,
		ButtonTopLeft:     buttonTopLeft,
		ButtonTopRight:    buttonTopRight,
		ButtonBottomLeft:  buttonBottomLeft,
		ButtonBottomRight: buttonBottomRight,
	}
}

func NewRemoteControl(
	radio models.Radio,
	roll, pitch, yaw, throttle joystick,
	btnFrontLeft, btnFrontRight button,
	btnTopLeft, btnTopRight button,
	btnBottomLeft, btnBottomRight button,
	commandPerSecond int,
	buzzer *piezobuzzer.Buzzer,
) *remoteControl {
	return &remoteControl{
		radio:            radio,
		roll:             roll,
		pitch:            pitch,
		yaw:              yaw,
		throttle:         throttle,
		btnFrontLeft:     btnFrontLeft,
		btnFrontRight:    btnFrontRight,
		btnTopLeft:       btnTopLeft,
		btnTopRight:      btnTopRight,
		btnBottomLeft:    btnBottomLeft,
		btnBottomRight:   btnBottomRight,
		commandPerSecond: commandPerSecond,
		buzzer:           buzzer,
	}
}

func (rc *remoteControl) Start(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer log.Println("Remote Control is stopped.")

		// log.Println("Waiting for connection...")
		var commandInterval time.Duration = time.Second / time.Duration(rc.commandPerSecond)

		var receiverChanOpen bool = true
		var connectionChanOpen bool = true
		var transmitterOpen bool = true
		var lastReading time.Time = time.Now()

		for transmitterOpen || connectionChanOpen || receiverChanOpen {
			select {
			case <-ctx.Done():
				if transmitterOpen {
					rc.radio.CloseTransmitter()
					transmitterOpen = false
				}

			case _, ok := <-rc.radio.GetReceiver():
				receiverChanOpen = ok

			case connectionState, ok := <-rc.radio.GetConnection():
				if ok {
					rc.setRadioConnectionState(connectionState)
				}
				connectionChanOpen = ok

			default:
			}

			if transmitterOpen && time.Since(lastReading) >= commandInterval {
				lastReading = time.Now()
				fc := rc.read()
				rc.radio.Transmit(fc)
				rc.shutdownPressed(fc.ButtonBottomRight, cancel)
				rc.suppressLostConnectionPressed(fc.ButtonBottomLeft)
			}
		}
	}()
}

func (rc *remoteControl) setRadioConnectionState(connectionState radio.ConnectionState) {
	rc.connectionState = connectionState
	switch rc.connectionState {
	case radio.CONNECTED:
		log.Println("Connected to Drone.")
		rc.buzzer.Stop()
		rc.buzzer.PlayNotes(piezobuzzer.ConnectedSound)
	case radio.DISCONNECTED:
		log.Println("Waiting for connection.")
		rc.buzzer.Stop()
		rc.buzzer.PlayNotes(piezobuzzer.DisconnectedSound)
	case radio.LOST:
		log.Println("Connection is lost.")
		rc.buzzer.WaveGenerator(piezobuzzer.WarningSound)
	}
}
