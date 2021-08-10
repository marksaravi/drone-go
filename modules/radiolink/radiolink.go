package radiolink

type FlightData struct {
	Roll          float32
	Pitch         float32
	Yaw           float32
	Throttle      float32
	Altitude      float32
	MotorsEngaged bool
}

type RadioLink interface {
	ReceiverOn()
	TransmitterOn()
	IsPayloadAvailable() bool
	TransmitFlightData(FlightData) error
	ReceiveFlightData() FlightData
}
