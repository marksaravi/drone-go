package connectors

type SPIConfig struct {
	BusNumber  int `yaml:"bus_number"`
	ChipSelect int `yaml:"chip_select"`
}
