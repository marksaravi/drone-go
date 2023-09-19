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

type remote struct {
	transmitter radioTransmiter

	lastCommand      time.Time
	commandPerSecond int
}

type RemoteCongigs struct {
	Transmitter radioTransmiter
}

func NewRemote(configs RemoteCongigs) *remote {
	return &remote{
		transmitter: configs.Transmitter,
		lastCommand: time.Now(),
	}
}

func (r *remote) Start(ctx context.Context) {
	running := true

	for running {
		select {
		default:
			if ok, commands := r.ReadCommands(); ok {
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
