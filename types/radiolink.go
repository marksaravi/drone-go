package types

type RadioLinkGPIOPins struct {
	CE string `yaml:"ce"`
}

type RadioLinkConfig struct {
	GPIO       RadioLinkGPIOPins `yaml:"gpio"`
	RxAddress  string            `yaml:"rx_address"`
	PowerDBm   byte              `yaml:"power_dBm"`
	BusNumber  int               `yaml:"bus_number"`
	ChipSelect int               `yaml:"chip_select"`
}

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
