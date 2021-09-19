package radioreceiver

import (
	"context"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radio interface {
	ReceiverOn()
	Receive() ([]byte, bool)
	TransmitterOn()
	Transmit([]byte) error
}
type radioReceiver struct {
	command    chan models.FlightCommands
	connection chan bool
}

func NewRadioReceiver(
	ctx context.Context,
	radio radio,
	commandPerSecond int,
	heartBeatPerSecond int,
	timeout time.Duration,
) *radioReceiver {
	commandChan := make(chan models.FlightCommands)
	connectionChan := make(chan bool)

	go receiverTask(ctx, radio, commandPerSecond, heartBeatPerSecond, timeout, commandChan, connectionChan)

	return &radioReceiver{
		command:    commandChan,
		connection: connectionChan,
	}
}

func receiverTask(
	ctx context.Context,
	radio radio,
	commandPerSecond int,
	heartBeatPerSecond int,
	timeout time.Duration,
	command chan models.FlightCommands,
	connection chan bool,
) {
	radio.ReceiverOn()
	receiveTicker := utils.NewTicker(ctx, commandPerSecond, 0)
	heartBeatTicker := utils.NewTicker(ctx, heartBeatPerSecond, 0)
	connected := false
	var lastDataTime time.Time
	for {
		select {
		case <-ctx.Done():
			return
		case <-receiveTicker:
			data, dataAvailable := radio.Receive()
			if dataAvailable {
				if !connected {
					connected = true
					connection <- true
				}
				lastDataTime = time.Now()
				command <- utils.DeserializeFlightCommand(data)
			} else {
				if connected && time.Since(lastDataTime) > timeout {
					connected = false
					connection <- false
				}
			}
		case <-heartBeatTicker:
			radio.TransmitterOn()
			radio.ReceiverOn()
		}
	}
}
