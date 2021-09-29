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
	var id uint32 = 0
	transmitter := radiotransmitter.NewRadioTransmitter(ctx, rc.radio, 20, 4, time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-transmitter.DataReadTicker:
			flightcommands := rc.read()
			flightcommands.Time = t
			flightcommands.Id = id
			id++
			transmitter.FlightComands <- flightcommands
			// fmt.Println(time.Unix(0, t).Format("2006-01-02 15:04:05.000"), flightcommands)
		case hb := <-transmitter.DroneHeartBeat:
			if hb {
				log.Println("Connected to Drone")
			} else {
				log.Println("Lost connection to Drone")
			}
		}
	}
}
