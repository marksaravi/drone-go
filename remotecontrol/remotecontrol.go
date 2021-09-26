package remotecontrol

import (
	"context"
	"log"
	"time"

	"github.com/marksaravi/drone-go/devices/radiotransmitter"
	"github.com/marksaravi/drone-go/models"
)

type button interface {
	Read() bool
}

type joystick interface {
	Read() float32
}

type remoteControl struct {
	radio        models.RadioLink
	roll         joystick
	pitch        joystick
	yaw          joystick
	throttle     joystick
	btnFrontLeft button
}

func (rc *remoteControl) read() models.FlightCommands {
	return models.FlightCommands{
		Roll:            rc.roll.Read(),
		Pitch:           rc.pitch.Read(),
		Yaw:             rc.yaw.Read(),
		Throttle:        rc.throttle.Read(),
		ButtonFrontLeft: rc.btnFrontLeft.Read(),
	}
}

func NewRemoteControl(radio models.RadioLink, roll, pitch, yaw, throttle joystick, btnFrontLeft button) *remoteControl {
	return &remoteControl{
		radio:        radio,
		roll:         roll,
		pitch:        pitch,
		yaw:          yaw,
		throttle:     throttle,
		btnFrontLeft: btnFrontLeft,
	}
}

func (rc *remoteControl) Start(ctx context.Context) {
	heartbeat := radiotransmitter.NewRadioTransmitter(ctx, rc.read, rc.radio, 40, 4, time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case hb := <-heartbeat:
			if hb {
				log.Println("Connected to Drone")
			} else {
				log.Println("Lost connection to Drone")
			}
		}
	}
}
