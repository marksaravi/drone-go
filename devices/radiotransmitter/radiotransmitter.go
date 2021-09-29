package radiotransmitter

import (
	"context"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radioTransmitter struct {
	FlightComands  chan<- models.FlightCommands
	DataReadTicker <-chan int64
	DroneHeartBeat <-chan bool
}

func NewRadioTransmitter(
	ctx context.Context,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
	hearbeattTimeout time.Duration,
) *radioTransmitter {
	heartbeatChan := make(chan bool)
	dataReadTicker := utils.NewTicker(ctx, commandPerSecond, 0)
	flightCommandsChan := make(chan models.FlightCommands)
	go transmitterRoutine(
		ctx,
		flightCommandsChan,
		heartbeatChan,
		radio,
		commandPerSecond,
		heartBeatPerSecond,
		hearbeattTimeout,
	)
	return &radioTransmitter{
		FlightComands:  flightCommandsChan,
		DataReadTicker: dataReadTicker,
		DroneHeartBeat: heartbeatChan,
	}
}

func transmitterRoutine(
	ctx context.Context,
	flightcommands chan models.FlightCommands,
	heartbeat chan bool,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
	hearbeattTimeout time.Duration,
) {
	heartbeatTicker := utils.NewTicker(ctx, int(float32(heartBeatPerSecond)*1.5), 0)
	var id uint32 = 0
	lastHeartbeat := time.Now()
	var heartbeating bool = false
	radio.ReceiverOn()
	for {
		select {
		case <-ctx.Done():
			return
		case fc := <-flightcommands:
			radio.TransmitterOn()
			radio.Transmit(utils.SerializeFlightCommand(fc))
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
