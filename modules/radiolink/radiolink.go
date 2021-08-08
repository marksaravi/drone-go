package radiolink

import "github.com/MarkSaravi/drone-go/types"

type RadioLink interface {
	ReceiverOn()
	TransmitterOn()
	IsPayloadAvailable() bool
	SendPayload(types.FlightData)
	ReadPayload() types.FlightData
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

func (r *radio) SendPayload(payload types.FlightData) {

}

func (r *radio) ReadPayload() types.FlightData {
	return types.FlightData{}
}
