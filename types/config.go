package types

import "github.com/MarkSaravi/drone-go/hardware/icm20948"

type PCA9685Config struct {
	Device          string        `yaml:"device"`
	PowerBrokerGPIO string        `yaml:"power_breaker_gpio"`
	Motors          map[int]Motor `yaml:"motors"`
}

type NRF204Config struct {
	BusNumber   int    `yaml:"bus_number"`
	ChipSelect  int    `yaml:"chip_select"`
	CEGPIO      string `yaml:"ce_gpio"`
	RxTxAddress string `yaml:"rx_tx_address"`
	PowerDBm    string `yaml:"power_dbm"`
}

type ImuConfig struct {
	ImuDataPerSecond            int     `yaml:"imu_data_per_second"`
	AccLowPassFilterCoefficient float64 `yaml:"acc_lowpass_filter_coefficient"`
	LowPassFilterCoefficient    float64 `yaml:"lowpass_filter_coefficient"`
}

type PidConfig struct {
	ProportionalGain float32 `yaml:"proportional–gain"`
	IntegralGain     float32 `yaml:"integral-gain"`
	DerivativeGain   float32 `yaml:"derivative-gain"`
}

type Motor struct {
	Label      string `yaml:"label"`
	ESCChannel int    `yaml:"esc_channel"`
}

type EscConfig struct {
	UpdateFrequency int     `yaml:"update_frequency"`
	MaxThrottle     float32 `yaml:"max_throttle"`
}

type FlightConfig struct {
	PID PidConfig `yaml:"pid"`
	Imu ImuConfig `yaml:"imu"`
	Esc EscConfig `yaml:"esc"`
}

type HardwareConfig struct {
	ICM20948 icm20948.ICM20948Config `yaml:"icm20948"`
	PCA9685  PCA9685Config           `yaml:"pca9685"`
	NRF204   NRF204Config            `yaml:"nrf204"`
}

type UdpLoggerConfig struct {
	Enabled          bool   `yaml:"enabled"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	PacketsPerSecond int    `yaml:"packets_per_second"`
	MaxDataPerPacket int    `yaml:"max_data_per_packet"`
}

type ApplicationConfig struct {
	Flight   FlightConfig    `yaml:"flight_control"`
	Hardware HardwareConfig  `yaml:"devices"`
	UDP      UdpLoggerConfig `yaml:"udp"`
}
