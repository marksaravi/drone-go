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
	Read() float32
}

var id uint32 = 0

type remoteControl struct {
	commandPerSecond  int
	radio             models.Radio
	roll              joystick
	pitch             joystick
	yaw               joystick
	throttle          joystick
	btnFrontLeft      button
	btnFrontRight     button
	btnTopLeft        button
	btnTopRight       button
	btnBottomLeft     button
	btnBottomRight    button
	buzzer            *piezobuzzer.Buzzer
	connectionState   radio.ConnectionState
	shutdownCountdown time.Time
}

func (rc *remoteControl) read() models.FlightCommands {
	id++
	return models.FlightCommands{
		Id:                id,
		Time:              time.Now().UnixNano(),
		Roll:              rc.roll.Read(),
		Pitch:             rc.pitch.Read(),
		Yaw:               rc.yaw.Read(),
		Throttle:          rc.throttle.Read(),
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

		log.Println("Waiting for connection...")
		var commandInterval time.Duration = time.Second / time.Duration(rc.commandPerSecond)
		var lastReading time.Time = time.Now()

		var connectionState radio.ConnectionState
		var connectionChanOpen bool = true
		var flightCommands models.FlightCommands
		var receiverChanOpen bool = true
		var running bool = true

		for running || connectionChanOpen || receiverChanOpen {
			if running && time.Since(lastReading) >= commandInterval {
				lastReading = time.Now()
				fc := rc.read()
				rc.radio.Transmit(fc)
				rc.processShutdownPressed(fc.ButtonBottomRight, cancel)
			}

			select {
			case flightCommands, receiverChanOpen = <-rc.radio.GetReceiver():
				if receiverChanOpen {
					log.Println("flight command: ", flightCommands.PayloadType)
				}
			default:
			}

			select {
			case connectionState, connectionChanOpen = <-rc.radio.GetConnection():
				if connectionChanOpen {
					rc.setRadioConnectionState(connectionState)
				}
			default:
			}

			select {
			case <-ctx.Done():
				if running {
					rc.radio.CloseTransmitter()
					running = false
				}
			default:
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
