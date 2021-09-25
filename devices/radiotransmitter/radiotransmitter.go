package radiotransmitter

import (
	"context"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

func NewRadioTransmitter(
	ctx context.Context,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
	hearbeattTimeout time.Duration,
) (chan models.FlightCommands, chan bool) {
	commandChan := make(chan models.FlightCommands)
	connectionChan := make(chan bool)

	go transmitterRoutine(ctx, commandChan, connectionChan, radio, commandPerSecond, heartBeatPerSecond, hearbeattTimeout)
	return commandChan, connectionChan
}

func transmitterRoutine(ctx context.Context, command chan models.FlightCommands, connection chan bool, radio models.RadioLink, commandPerSecond int, heartBeatPerSecond int, hearbeattTimeout time.Duration) {
	transmitterTicker := utils.NewTicker(ctx, commandPerSecond, 0)
	heartbeatTicker := utils.NewTicker(ctx, int(float32(heartBeatPerSecond)*1.5), 0)
	var id uint32 = 0
	lastHeartbeat := time.Now()
	var connected bool = false
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-transmitterTicker:
			radio.TransmitterOn()
			radio.Transmit(utils.SerializeFlightCommand(models.FlightCommands{
				Id:   id,
				Time: t,
			}))
			id++
			radio.ReceiverOn()
		case <-heartbeatTicker:
			if _, available := radio.Receive(); available {
				lastHeartbeat = time.Now()
				if !connected {
					connected = true
					connection <- true
				}
			} else {
				if time.Since(lastHeartbeat) > hearbeattTimeout {
					if connected {
						connected = false
						connection <- false
					}
				}
			}
		}
	}
}
