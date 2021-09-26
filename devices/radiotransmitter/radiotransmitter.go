package radiotransmitter

import (
	"context"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type commandReader interface {
	Read() models.FlightCommands
}

func NewRadioTransmitter(
	ctx context.Context,
	commandreader func() models.FlightCommands,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
	hearbeattTimeout time.Duration,
) chan bool {
	heartbeatChan := make(chan bool)

	go transmitterRoutine(
		ctx,
		commandreader,
		heartbeatChan,
		radio,
		commandPerSecond,
		heartBeatPerSecond,
		hearbeattTimeout,
	)
	return heartbeatChan
}

func transmitterRoutine(
	ctx context.Context,
	commandreader func() models.FlightCommands,
	heartbeat chan bool,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
	hearbeattTimeout time.Duration,
) {
	transmitterTicker := utils.NewTicker(ctx, commandPerSecond, 0)
	heartbeatTicker := utils.NewTicker(ctx, int(float32(heartBeatPerSecond)*1.5), 0)
	var id uint32 = 0
	lastHeartbeat := time.Now()
	var heartbeating bool = false
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-transmitterTicker:
			radio.TransmitterOn()
			flightcommands := commandreader()
			flightcommands.Time = t
			radio.Transmit(utils.SerializeFlightCommand(flightcommands))
			id++
			radio.ReceiverOn()
		case <-heartbeatTicker:
			if _, available := radio.Receive(); available {
				lastHeartbeat = time.Now()
				if !heartbeating {
					heartbeating = true
					heartbeat <- true
				}
			} else {
				if time.Since(lastHeartbeat) > hearbeattTimeout {
					if heartbeating {
						heartbeating = false
						heartbeat <- false
					}
				}
			}
		}
	}
}
