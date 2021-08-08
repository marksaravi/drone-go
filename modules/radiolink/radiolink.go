package radiolink

import "github.com/MarkSaravi/drone-go/types"

type RadioLink interface {
	ReceiverOn()
	TransmitterOn()
	IsPayloadAvailable() bool
	TransmitFlightData(types.FlightData)
	ReceiveFlightData() types.FlightData
}

type radio struct {
}

func NewRadioLink() RadioLink {
	return &radio{}
}

func (r *radio) ReceiverOn() {

}

func (r *radio) TransmitterOn() {

}

func (r *radio) IsPayloadAvailable() bool {
	return false
}

func (r *radio) TransmitFlightData(payload types.FlightData) {

}

func (r *radio) ReceiveFlightData() types.FlightData {
	return types.FlightData{}
}
