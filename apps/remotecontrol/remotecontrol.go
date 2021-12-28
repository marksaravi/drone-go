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
	Read() uint16
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
	return models.FlightCommands{
		Roll:              uint8(rc.roll.Read()),
		Pitch:             uint8(rc.pitch.Read()),
		Yaw:               uint8(rc.yaw.Read()),
		Throttle:          uint8(rc.throttle.Read()),
		ButtonFrontLeft:   rc.btnFrontLeft.Read(),
		ButtonFrontRight:  rc.btnFrontRight.Read(),
		ButtonTopLeft:     rc.btnTopLeft.Read(),
		ButtonTopRight:    rc.btnTopRight.Read(),
		ButtonBottomLeft:  rc.btnBottomLeft.Read(),
		ButtonBottomRight: rc.btnBottomRight.Read(),
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
