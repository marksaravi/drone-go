package types

type RadioLinkGPIOPins struct {
	CE string `yaml:"ce"`
}

type RadioLinkConfig struct {
	GPIO       RadioLinkGPIOPins `yaml:"gpio"`
	RxAddress  string            `yaml:"rx_address"`
	BusNumber  int               `yaml:"bus_number"`
	ChipSelect int               `yaml:"chip_select"`
}
