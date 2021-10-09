package remotecontrol

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices/radiotransmitter"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
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

func (rc *remoteControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	var id uint32 = 0
	configs := config.ReadRemoteControlConfig().Radio
	log.Println(configs.CommandPerSecond)
	dataReadTicker := utils.NewTicker(ctx, "Remote Control", configs.CommandPerSecond)
	command, connection := radiotransmitter.NewRadioTransmitter(ctx, wg)
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-dataReadTicker:
			flightcommands := rc.read()
			flightcommands.Time = t
			flightcommands.Id = id
			id++
			command <- flightcommands
		case connected := <-connection:
			if connected {
				log.Println("Connected to Drone")
			} else {
				log.Println("Lost connection to Drone")
			}
		}
	}
}
