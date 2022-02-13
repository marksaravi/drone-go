package radio

import (
	"github.com/marksaravi/drone-go/models"
)

type ConnectionState = int

const (
	IDLE ConnectionState = iota
	DISCONNECTED
	CONNECTED
	CONNECTION_LOST
)

const (
	NO_COMMAND models.FlightCommandType = iota
	COMMAND
	HEARTBEAT
	RECEIVER_OFF
)

type radioLink interface {
	PowerOn()
	PowerOff()
	ClearStatus()
}
type radioTransmitterLink interface {
	radioLink
	TransmitterOn()
	Transmit(models.Payload) error
	IsTransmitFailed(update bool) bool
}

type radioReceiverLink interface {
	radioLink
	ReceiverOn()
	Listen()
	Receive() (models.Payload, error)
	IsReceiverDataReady(update bool) bool
}
type radioReceiver struct {
	receiveChannel    chan models.FlightCommands
	connectionChannel chan ConnectionState
	radiolink         radioReceiverLink
	connectionState   ConnectionState
}

type radioTransmitter struct {
	TransmitChannel   chan models.FlightCommands
	connectionChannel chan ConnectionState
	radiolink         radioTransmitterLink
	connectionState   ConnectionState
}
