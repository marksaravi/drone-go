package remotecontrol

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/constants"
	piezobuzzer "github.com/marksaravi/drone-go/hardware/piezo-buzzer"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type button interface {
	Read() bool
}

type oledDisplay interface {
	Println(msg string, y int)
}

type joystick interface {
	Read() int
}

type radioTransmitter interface {
	Transmit(models.FlightCommands)
	GetConnectionStateChannel() <-chan int
	Close()
	SuppressLostConnection()
}
type remoteControl struct {
	commandPerSecond  int
	radio             radioTransmitter
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
	display           oledDisplay
	buzzer            *piezobuzzer.Buzzer
	shutdownCountdown time.Time
}

var lastPrinted time.Time = time.Now()

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

	utils.PrintIntervally(fmt.Sprintf("%6d, %6d, %6d, %6d, %v, %v\n", roll, pitch, yaw, throttle, buttonTopLeft, buttonTopRight), "remotedata", time.Second, false)

	return models.FlightCommands{
		Roll:              uint16(roll),
		Pitch:             uint16(pitch),
		Yaw:               uint16(yaw),
		Throttle:          uint16(throttle),
		ButtonFrontLeft:   buttonFrontLeft,
		ButtonFrontRight:  buttonFrontRight,
		ButtonTopLeft:     buttonTopLeft,     // used for calibration -
		ButtonTopRight:    buttonTopRight,    // used for calibration +
		ButtonBottomLeft:  buttonBottomLeft,  // used for stopping disconnect alarm
		ButtonBottomRight: buttonBottomRight, // used for shutting down the remote control
	}
}

func NewRemoteControl(
	radio radioTransmitter,
	roll, pitch, yaw, throttle joystick,
	btnFrontLeft, btnFrontRight button,
	btnTopLeft, btnTopRight button,
	btnBottomLeft, btnBottomRight button,
	commandPerSecond int,
	display oledDisplay,
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
		display:          display,
		buzzer:           buzzer,
	}
}

func (rc *remoteControl) Start(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer time.Sleep(time.Second)
		defer rc.buzzer.PlaySound(piezobuzzer.ExitSound)
		defer log.Println("Remote Control is stopped.")

		var readingInterval time.Duration = time.Second / time.Duration(rc.commandPerSecond)
		var connectionChanOpen bool = true
		var transmitterOpen bool = true
		var lastReading time.Time = time.Now()

		for transmitterOpen || connectionChanOpen {
			select {
			case <-ctx.Done():
				if transmitterOpen {
					transmitterOpen = false
					rc.radio.Close()
				}

			case connectionState, ok := <-rc.radio.GetConnectionStateChannel():
				if ok {
					rc.setRadioConnectionState(connectionState)
				}
				connectionChanOpen = ok

			default:
				if transmitterOpen && time.Since(lastReading) >= readingInterval {
					lastReading = time.Now()
					fc := rc.read()
					rc.radio.Transmit(fc)
					rc.actOnShutdownButtonState(fc.ButtonBottomRight, cancel)
					rc.actOnSuppressLostConnectionButtonState(fc.ButtonBottomLeft)
				}
			}
		}
	}()
}

func (rc *remoteControl) setRadioConnectionState(connectionState int) {
	switch connectionState {
	case constants.CONNECTED:
		log.Println("Connected to Drone.")
		rc.display.Println("Connected!", 3)
		rc.buzzer.PlaySound(piezobuzzer.ConnectedSound)
	case constants.WAITING_FOR_CONNECTION:
		log.Println("Waiting for connection.")
		rc.display.Println("Waiting...", 3)
		rc.buzzer.PlaySound(piezobuzzer.DisconnectedSound)
	case constants.DISCONNECTED:
		log.Println("Connection is lost.")
		rc.display.Println("Drone is Lost", 3)
		// rc.buzzer.WaveGenerator(piezobuzzer.WarningSound)
		rc.buzzer.PlaySound(piezobuzzer.DisconnectedSound)
	}
}
