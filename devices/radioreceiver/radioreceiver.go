package radioreceiver

import (
	"context"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radioReceiver struct {
	command    chan models.FlightCommands
	connection chan bool
}

func NewRadioReceiver(
	ctx context.Context,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
	receiverTimeout time.Duration,
) *radioReceiver {
	commandChan := make(chan models.FlightCommands)
	connectionChan := make(chan bool)

	go receiverRoutine(ctx, radio, int(float32(commandPerSecond)*1.5), heartBeatPerSecond, receiverTimeout, commandChan, connectionChan)

	return &radioReceiver{
		command:    commandChan,
		connection: connectionChan,
	}
}

func receiverRoutine(
	ctx context.Context,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
	receiverTimeout time.Duration,
	command chan models.FlightCommands,
	connection chan bool,
) {
	radio.ReceiverOn()
	receiveTicker := utils.NewTicker(ctx, commandPerSecond, 0)
	heartBeatTicker := utils.NewTicker(ctx, heartBeatPerSecond, 0)
	connected := false
	var lastDataTime time.Time = time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case <-receiveTicker:
			data, dataAvailable := radio.Receive()
			if dataAvailable {
				lastDataTime = time.Now()
				if !connected {
					connected = true
					connection <- true
				}
				command <- utils.DeserializeFlightCommand(data)
			} else {
				if connected && time.Since(lastDataTime) > receiverTimeout {
					connected = false
					connection <- false
				}
			}
		case <-heartBeatTicker:
			radio.TransmitterOn()
			timedata := utils.Int64ToBytes(time.Now().UnixNano())
			radio.Transmit(utils.SliceToArray32(timedata[:]))
			radio.ReceiverOn()
		}
	}
}
