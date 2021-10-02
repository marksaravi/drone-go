package radioreceiver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/drivers/nrf204"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radioReceiver struct {
	Command    chan models.FlightCommands
	Connection chan bool
}

func NewRadioReceiver(ctx context.Context, wg *sync.WaitGroup) (<-chan models.FlightCommands, <-chan bool) {
	radio := nrf204.NewRadio()
	config := config.ReadFlightControlConfig().Radio
	receiver := newRadioReceiver(ctx, wg, radio, config.CommandPerSecond, config.CommandTimeoutMS, config.HeartBeatPerSecond)
	return receiver.Command, receiver.Connection
}

func newRadioReceiver(
	ctx context.Context,
	wg *sync.WaitGroup,
	radio models.RadioLink,
	commandPerSecond int,
	commandTimeoutMs int,
	heartBeatPerSecond int,
) *radioReceiver {
	commandChan := make(chan models.FlightCommands)
	connectionChan := make(chan bool)

	wg.Add(1)
	go receiverRoutine(ctx, wg, radio, commandPerSecond, commandTimeoutMs, heartBeatPerSecond, commandChan, connectionChan)

	return &radioReceiver{
		Command:    commandChan,
		Connection: connectionChan,
	}
}

func receiverRoutine(
	ctx context.Context,
	wg *sync.WaitGroup,
	radio models.RadioLink,
	commandPerSecond int,
	commandTimeoutMs int,
	heartBeatPerSecond int,
	command chan models.FlightCommands,
	connection chan bool,
) {
	defer wg.Done()
	defer log.Println("RADIO CLOSED")
	defer close(command)
	defer close(connection)

	radio.ReceiverOn()
	receiveTicker := utils.NewTicker(ctx, wg, commandPerSecond*2, 0)
	commandTimeout := time.Millisecond * time.Duration(commandTimeoutMs)
	heartBeatTicker := utils.NewTicker(ctx, wg, heartBeatPerSecond, 0)
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
				if connected && time.Since(lastDataTime) > commandTimeout {
					connected = false
					connection <- false
				}
			}
		case <-heartBeatTicker:
			radio.TransmitterOn()
			timedata := utils.Int64ToBytes(time.Now().UnixNano())
			if err := radio.Transmit(utils.SliceToArray32(timedata[:])); err != nil {
				fmt.Println(err)
			}
			radio.ReceiverOn()
		}
	}
}
