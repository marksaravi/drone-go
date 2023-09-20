package remote

import (
	"context"
	"time"
)

type radioTransmiter interface {
	On()
	Transmit(payload []byte) error
}

type joystick interface {
	Read() uint16
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

	lastCommand      time.Time
	commandPerSecond int
}

type RemoteSettings struct {
	Transmitter                radioTransmiter
	CommandPerSecond           int
	Roll, Pitch, Yaw, Throttle joystick
}

func NewRemoteControl(settings RemoteSettings) *remoteControl {
	return &remoteControl{
		transmitter:      settings.Transmitter,
		commandPerSecond: settings.CommandPerSecond,
		roll:             settings.Roll,
		pitch:            settings.Pitch,
		yaw:              settings.Yaw,
		throttle:         settings.Throttle,
		lastCommand:      time.Now(),
	}
}

func (r *remoteControl) Start(ctx context.Context) {
	running := true
	r.transmitter.On()
	for running {
		select {
		default:
			if commands, ok := r.ReadCommands(); ok {
				// log.Println(commands)
				r.transmitter.Transmit([]byte{
					byte(commands.roll),
					byte(commands.pitch),
					byte(commands.yaw),
					byte(commands.throttle),
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				})

			}
		case <-ctx.Done():
			running = false
		}
	}

}
