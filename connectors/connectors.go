package connectors

import "periph.io/x/periph/conn/spi"

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
