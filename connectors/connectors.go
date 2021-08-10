package connectors

import "periph.io/x/periph/conn/spi"

type SPIConfig struct {
	BusNumber   int `yaml:"bus_number"`
	ChipSelect  int `yaml:"chip_select"`
	Mode        int `yaml:"mode"`
	SpeedMegaHz int `yaml:"speed-mega-hz"`
}

func ConfigToSPIMode(configValue int) spi.Mode {
	switch configValue {
	case 0:
		return spi.Mode0
	case 1:
		return spi.Mode1
	case 2:
		return spi.Mode2
	case 3:
		return spi.Mode3
	default:
		return spi.Mode0
	}
}
