package radiotransmitter

import (
	"context"
	"sync"
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
	wg *sync.WaitGroup,
	radio models.RadioLink,
	commandPerSecond int,
	hearbeattTimeout time.Duration,
) *radioTransmitter {
	heartbeatChan := make(chan bool)
	dataReadTicker := utils.NewTicker(ctx, commandPerSecond, 0)
	flightCommandsChan := make(chan models.FlightCommands)
	wg.Add(1)
	go transmitterRoutine(
		ctx,
		wg,
		flightCommandsChan,
		heartbeatChan,
		radio,
		commandPerSecond,
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
	wg *sync.WaitGroup,
	flightcommands chan models.FlightCommands,
	heartbeat chan bool,
	radio models.RadioLink,
	commandPerSecond int,
	hearbeattTimeout time.Duration,
) {
	defer wg.Done()
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
		default:
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
