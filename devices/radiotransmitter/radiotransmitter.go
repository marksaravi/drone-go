package radiotransmitter

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/drivers/nrf204"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radioTransmitter struct {
	Command    chan<- models.FlightCommands
	Connection <-chan bool
}

func NewRadioTransmitter(
	ctx context.Context,
	wg *sync.WaitGroup,
) (
	chan<- models.FlightCommands,
	<-chan bool,

) {
	radio := nrf204.NewRadio()
	configs := config.ReadRemoteControlConfig().Radio
	transmitter := newRadioTransmitter(ctx, wg, radio, configs.HeartBeatTimeoutMS)
	return transmitter.Command, transmitter.Connection
}

func newRadioTransmitter(
	ctx context.Context,
	wg *sync.WaitGroup,
	radio models.RadioLink,
	hearbeatTimeoutMs int,

) *radioTransmitter {
	command := make(chan models.FlightCommands, 2)
	connection := make(chan bool, 2)
	transmitterRoutine(ctx, wg, radio, command, connection, hearbeatTimeoutMs)
	return &radioTransmitter{
		Command:    command,
		Connection: connection,
	}
}

func transmitterRoutine(
	ctx context.Context,
	wg *sync.WaitGroup,
	radio models.RadioLink,
	command chan models.FlightCommands,
	connection chan bool,
	heartbeatTimeoutMs int,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer log.Println("RADIO CLOSED")
		radio.ReceiverOn()
		heartbeatTimeout := time.Duration(heartbeatTimeoutMs) * time.Millisecond
		log.Print(heartbeatTimeout)
		var connected bool = false
		var lastHeartbeat time.Time
		for {
			select {
			case flightcommands := <-command:
				_, isAvailable := radio.Receive()
				if isAvailable && connected {
					lastHeartbeat = time.Now()
				} else if isAvailable && !connected {
					lastHeartbeat = time.Now()
					connected = true
					connection <- true
				} else if !isAvailable && connected {
					if time.Since(lastHeartbeat) > heartbeatTimeout {
						connected = false
						connection <- false
					}
				}
				radio.TransmitterOn()
				radio.Transmit(utils.SerializeFlightCommand(flightcommands))
				radio.ReceiverOn()
			case <-ctx.Done():
				return
			}
		}
	}()
}
