package types

type Offsets struct {
	X int16 `yaml:"X"`
	Y int16 `yaml:"Y"`
	Z int16 `yaml:"Z"`
}

// AccelerometerConfig is the configurations for Accelerometer
type AccelerometerConfig struct {
	SensitivityLevel     string  `yaml:"sensitivity_level"`
	LowPassFilterEnabled bool    `yaml:"lowpass_filter_enabled"`
	LowPassFilterConfig  int     `yaml:"lowpass_filter_config"`
	Averaging            int     `yaml:"averaging"`
	Offsets              Offsets `yaml:"offsets"`
}

// GyroscopeConfig is the configuration for Gyroscope
type GyroscopeConfig struct {
	SensitivityLevel     string  `yaml:"sensitivity_level"`
	LowPassFilterEnabled bool    `yaml:"lowpass_filter_enabled"`
	LowPassFilterConfig  int     `yaml:"lowpass_filter_config"`
	Averaging            int     `yaml:"averaging"`
	Offsets              Offsets `yaml:"offsets"`
}

// MagnetometerConfig is the configuration for Magnetometer
type MagnetometerConfig struct {
	SensitivityLevel string `yaml:"sensitivity_level"`
}

type Icm20948Config struct {
	BusNumber  int                 `yaml:"bus_number"`
	ChipSelect int                 `yaml:"chip_select"`
	AccConfig  AccelerometerConfig `yaml:"accelerometer"`
	GyroConfig GyroscopeConfig     `yaml:"gyroscope"`
	MagConfig  MagnetometerConfig  `yaml:"magnetometer"`
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

type EscConfig struct {
	UpdateFrequency int     `yaml:"update_frequency"`
	MaxThrottle     float32 `yaml:"max_throttle"`
	SPIDevice       string  `yaml:"spi_device"`
	PowerBrokerGPIO string  `yaml:"power_breaker_gpio"`
}

type FlightConfig struct {
	PID PidConfig `yaml:"pid"`
	Imu ImuConfig `yaml:"imu"`
	Esc EscConfig `yaml:"esc"`
}

type HardwareConfig struct {
	ICM20948 Icm20948Config `yaml:"icm20948"`
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
