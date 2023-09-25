package remote

import (
	"context"
	"fmt"
	"time"
)

type radioTransmiter interface {
	On()
	Transmit(payload []byte) error
}

type joystick interface {
	Read() uint16
}

type PushButton interface {
	Name() string
	IsPressed() bool
}

type commands struct {
	roll     uint16
	pitch    uint16
	yaw      uint16
	throttle uint16
}

type remoteControl struct {
	transmitter radioTransmiter
	roll        joystick
	pitch       joystick
	yaw         joystick
	throttle    joystick
	buttons     []PushButton

	commandPerSecond int
}

type RemoteSettings struct {
	Transmitter                radioTransmiter
	CommandPerSecond           int
	Roll, Pitch, Yaw, Throttle joystick
	PushButtons                []PushButton
}

func NewRemoteControl(settings RemoteSettings) *remoteControl {
	return &remoteControl{
		transmitter:      settings.Transmitter,
		commandPerSecond: settings.CommandPerSecond,
		roll:             settings.Roll,
		pitch:            settings.Pitch,
		yaw:              settings.Yaw,
		throttle:         settings.Throttle,
		buttons:          settings.PushButtons,
	}
}

func (r *remoteControl) Start(ctx context.Context) {
	running := true
	r.transmitter.On()
	lastCommand:=time.Now()
	commandTimeout:=time.Second/time.Duration(r.commandPerSecond)
	for running {
		select {
		default:
			if time.Since(lastCommand)>=commandTimeout {
				if commands, ok := r.ReadCommands(); ok {
					payload:= []byte {
						byte(commands.roll),
						byte(commands.pitch),
						byte(commands.yaw),
						byte(commands.throttle),
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					}
					fmt.Println(payload)
					r.transmitter.Transmit(payload)
				}
				lastCommand=time.Now()
			}
		case <-ctx.Done():
			running = false
		}
	}
}
