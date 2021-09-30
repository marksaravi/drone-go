package radioreceiver

import (
	"context"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radioReceiver struct {
	Command    chan models.FlightCommands
	Connection chan bool
}

func NewRadioReceiver(
	ctx context.Context,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
) *radioReceiver {
	commandChan := make(chan models.FlightCommands)
	connectionChan := make(chan bool)

	go receiverRoutine(ctx, radio, commandPerSecond, heartBeatPerSecond, commandChan, connectionChan)

	return &radioReceiver{
		Command:    commandChan,
		Connection: connectionChan,
	}
}

func receiverRoutine(
	ctx context.Context,
	radio models.RadioLink,
	commandPerSecond int,
	heartBeatPerSecond int,
	command chan models.FlightCommands,
	connection chan bool,
) {
	radio.ReceiverOn()
	receiveTicker := utils.NewTicker(ctx, commandPerSecond*2, 0)
	receiverTimeout := time.Second / time.Duration(commandPerSecond/2)
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
			// log.Println("heartbeat")
			radio.TransmitterOn()
			timedata := utils.Int64ToBytes(time.Now().UnixNano())
			if err := radio.Transmit(utils.SliceToArray32(timedata[:])); err != nil {
				fmt.Println(err)
			}
			radio.ReceiverOn()
		}
	}
}
