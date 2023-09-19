package remote

import (
	"context"
	"log"
	"time"
)

type radioTransmiter interface {
	Start()
	Transmit(payload []byte) error
	PayloadSize() int
}

type commands struct {
	roll     byte
	pitch    byte
	yaw      byte
	throttle byte
}

type remoteControl struct {
	transmitter radioTransmiter

	lastCommand      time.Time
	commandPerSecond int
}

type RemoteCongigs struct {
	Transmitter      radioTransmiter
	CommandPerSecond int
}

func NewRemote(configs RemoteCongigs) *remoteControl {
	return &remoteControl{
		transmitter:      configs.Transmitter,
		commandPerSecond: configs.CommandPerSecond,
		lastCommand:      time.Now(),
	}
}

func (r *remoteControl) Start(ctx context.Context) {
	running := true

	for running {
		select {
		default:
			if commands, ok := r.ReadCommands(); ok {
				r.transmitter.Transmit([]byte{
					commands.roll,
					commands.pitch,
					commands.yaw,
					commands.throttle,
					0, 0, 0, 0,
				})
				log.Println(commands)
			}
		case <-ctx.Done():
			running = false
		}
	}

}
