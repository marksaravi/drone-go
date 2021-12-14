package remotecontrol

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/models"
)

type button interface {
	Read() bool
}

type joystick interface {
	Read() float32
}

type remoteControl struct {
	commandPerSecond int
	radio            models.Radio
	roll             joystick
	pitch            joystick
	yaw              joystick
	throttle         joystick
	btnFrontLeft     button
	btnFrontRight    button
	btnTopLeft       button
	btnTopRight      button
	btnBottomLeft    button
	btnBottomRight   button
}

func (rc *remoteControl) read() models.FlightCommands {
	return models.FlightCommands{
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
	}
}

func (rc *remoteControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer log.Println("Remote Control is stopped.")

		// var id uint32 = 0
		// dataReadTicker := utils.NewTicker(ctx, "Remote Control", rc.commandPerSecond)
		lastPrinted := time.Now()
		connection := rc.radio.GetConnection()
		log.Println("Waiting for connection...")
		for {
			select {
			case <-ctx.Done():
				rc.radio.Close()
				return
				// case t := <-dataReadTicker:
				// fc := rc.read()

			// 	fc.Time = t
			// 	fc.Id = id
			// 	id++
			// 	rc.radio.Transmit(fc)
			case connected := <-connection:
				if connected {
					log.Println("Connected to Drone")
				} else {
					log.Println("Lost connection to Drone")
				}
			default:
				fc := rc.read()
				if time.Since(lastPrinted) >= time.Second/4 {
					log.Println(fc)
					lastPrinted = time.Now()
				}
			}
		}
	}()
}
