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
	radio          models.RadioLink
	roll           joystick
	pitch          joystick
	yaw            joystick
	throttle       joystick
	btnFrontLeft   button
	btnFrontRight  button
	btnTopLeft     button
	btnTopRight    button
	btnBottomLeft  button
	btnBottomRight button
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
	radio models.RadioLink,
	roll, pitch, yaw, throttle joystick,
	btnFrontLeft, btnFrontRight button,
	btnTopLeft, btnTopRight button,
	btnBottomLeft, btnBottomRight button,
) *remoteControl {
	return &remoteControl{
		radio:          radio,
		roll:           roll,
		pitch:          pitch,
		yaw:            yaw,
		throttle:       throttle,
		btnFrontLeft:   btnFrontLeft,
		btnFrontRight:  btnFrontRight,
		btnTopLeft:     btnTopLeft,
		btnTopRight:    btnTopRight,
		btnBottomLeft:  btnBottomLeft,
		btnBottomRight: btnBottomRight,
	}
}

func (rc *remoteControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	var id uint32 = 0
	configs := config.ReadRemoteControlConfig().Radio
	dataReadTicker := utils.NewTicker(ctx, "Remote Control", configs.CommandPerSecond)
	command, connection := radiotransmitter.NewRadioTransmitter(ctx, wg)
	log.Println("Waiting for connection...")
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-dataReadTicker:
			fc := rc.read()

			fc.Time = t
			fc.Id = id
			id++
			command <- fc
		case connected := <-connection:
			if connected {
				log.Println("Connected to Drone")
			} else {
				log.Println("Lost connection to Drone")
			}
		}
	}
}
